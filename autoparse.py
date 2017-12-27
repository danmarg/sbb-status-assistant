# -*- coding: utf-8 -*-
import argparse
import copy
import collections
import csv
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
args = parser.parse_args()

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

print('Processed',len(entities),'items')

processed = collections.defaultdict(set)
for entity, values in entities.items():
    entity = re.sub(r'[()]', '', entity)  # Parens are not allowed in entity names.
    processed[entity] = processed[entity].union(values)

print('Merged into',len(processed),'items')

with open(args.output, 'w') as f:
    w = csv.writer(f, delimiter=',', quoting=csv.QUOTE_ALL)
    for entity, values in processed.items():
        w.writerow([entity] + list(values))

