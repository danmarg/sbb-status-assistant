# -*- coding: utf-8 -*-
# Parses data from https://data.sbb.ch/explore/ into proper Dialogflow
# entities.
import argparse
import copy
import collections
import csv
import json
import re

REPLACEMENTS = [
    (r'\s*Z[üu]rich\s*', ''),
    (r'\s*Bahnhof\s*', ''),
    (r'ü', 'ue'),
    (r'ü', 'u'),
    (r'ö', 'oe'),
    (r'ö', 'o'),
    (r'ä', 'ae'),
    (r'ä', 'a'),
    (r',', ''),
    (r'[()]', ''),
    (r'\s*\(.*\)\s*', ''),
]

parser = argparse.ArgumentParser()
parser.add_argument("input")
parser.add_argument("output")
parser.add_argument("json")
args = parser.parse_args()

def edit_distance(s1, s2):
    m=len(s1)+1
    n=len(s2)+1

    tbl = {}
    for i in range(m): tbl[i,0]=i
    for j in range(n): tbl[0,j]=j
    for i in range(1, m):
        for j in range(1, n):
            cost = 0 if s1[i-1] == s2[j-1] else 1
            tbl[i,j] = min(tbl[i, j-1]+1, tbl[i-1, j]+1, tbl[i-1, j-1]+cost)

    return tbl[i,j]

entities = {}
# The CSV has names in column 2 (0-indexed).
with open(args.input, 'r') as f:
    r = csv.reader(f, delimiter=';')
    for row in r:
      entities[row[2]] = set([row[2]])

for entity in entities.keys():
    unchanged = 0
    # Two passes to get multiple permutations.
    for i in range(len(REPLACEMENTS)*2):
        new = copy.copy(entities[entity])
        for e in entities[entity]:
            for r, v in REPLACEMENTS:
                x = re.sub(r, v, e, re.IGNORECASE)
                x = re.sub('(^[, ]+)|([, ]+$)', '', x)
                new.add(x)
        # quit early
        if len(new) == len(entities[entity]):
            unchanged += 1
        if unchanged > 2:
            break
        entities[entity] = new

reverse = {}  # Alternate values for entities, to try to ensure uniqueness.
for entity, alternates in entities.items():
    alternates = copy.copy(alternates)
    for alt in alternates:
        if alt in reverse:
            print('Found conflicting alternate',alt,'maps to',entity,'and',reverse[alt])
            # Prefer the version with lower edit distance from the original.
            if edit_distance(alt, reverse[alt]) < edit_distance(alt, entity):
                print('Preferring',reverse[alt])
                entities[entity].remove(alt)
            else:
                print('Preferring',entity)
                entities[reverse[alt]].remove(alt)
                reverse[alt] = entity
        else:
            reverse[alt] = entity

print('Processed',len(entities),'items')

processed = collections.defaultdict(set)
for entity, values in entities.items():
    entity = re.sub(r'[()]', '', entity)  # Parens are not allowed in entity names.
    processed[entity] = processed[entity].union(values)

print('Merged into',len(processed),'items')

with open(args.output, 'w') as f:
    if args.json:
        m = {}
        for entity, values in entities.items():
            for v in values:
                m[v.lower()] = entity
        f.write(json.dumps(m))
    else:
        w = csv.writer(f, delimiter=',', quoting=csv.QUOTE_ALL)
        for entity, values in processed.items():
            w.writerow([entity] + list([x for x in values
                                        if not '(' in x and not ')' in x]))
