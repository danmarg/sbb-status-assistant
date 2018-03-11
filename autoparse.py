# -*- coding: utf-8 -*-
# Parses data from https://data.sbb.ch/explore/dataset/station-didok/export/
# into proper Dialogflow entities.
import argparse
import copy
import collections
import csv
import distutils.util
import json
import re

REPLACEMENTS = [
    (r'\s*Z[üu]rich\s*', ''),
    (r'(^|[\s,/]+)Bahnhof([\s,/]+|$)', ''),
    (r'ü', 'ue'),
    (r'ü', 'u'),
    (r'ö', 'oe'),
    (r'ö', 'o'),
    (r'ä', 'ae'),
    (r'ä', 'a'),
    (r',', ''),
    (r',?\sHB', ' Hauptbahnhof'),
    (r',?\sHB', ''),
    (r',?\sHauptbahnhof', ' HB'),
    (r',?\sHauptbahnhof', ''),
    (r'[()]', ''),
    (r'\s*\(.*\)\s*', ''),
    (r'Basel', 'Bâle'),
    (r'Genèv', 'Genf'),
]

parser = argparse.ArgumentParser()
parser.add_argument('--input')
parser.add_argument('--output')
parser.add_argument('--json', dest='json', type=lambda x:bool(distutils.util.strtobool(x)))
parser.add_argument('--allowed_bus_operators')
parser.add_argument('--no_busses_except_for')
# Good default arg for the above is "Zürich,Basel,Bern,Luzern,Winterthur,Locarno,Lugano,Genèv,Laus,Gallen,Biel,Thun,Fribo,Köniz,Chaux,Schaffha,Vernier,Chur,Neuch,Uster"
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

allowed_busses = (set(args.allowed_bus_operators.split(','))
        if args.allowed_bus_operators else [])
no_busses_except_for = (set(args.no_busses_except_for.split(','))
        if args.no_busses_except_for else [])

entities = {}
# The CSV has names in column 2 (0-indexed).
with open(args.input, 'r') as f:
    r = csv.reader(f, delimiter=';')
    for row in r:
      if row[7] not in ('Haltestelle', 'Haltestelle_und_Bedienpunkt'):
        continue
      if allowed_busses and row[6] not in allowed_busses:
        continue
      # Filter on modes of transport, except for whitelisted cities.
      if no_busses_except_for and 'Zug' not in row[8]:
        if not [a for a in no_busses_except_for if a in row[2]]:
          continue
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
        if alt == "":
            continue
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
