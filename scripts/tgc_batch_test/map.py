import json

with open('tgc_resource_map.json') as f:
    mappings = json.load(f)

with open('resources_to_add.txt') as f:
    targets = [l.strip() for l in f if l.strip()]

matched_resources = {}
missing = []

for tgt in targets:
    found = False
    for product in mappings:
        for res in mappings[product]:
            if res['target'] == tgt:
                matched_resources[product] = matched_resources.get(product, [])
                matched_resources[product].append({'file': res['file'], 'target': tgt})
                found = True
                break
    if not found:
        missing.append(tgt)

if missing:
    print(f"Missing: {len(missing)}")
    for m in missing:
        print(m)

# Fix missing mappings based on previous troubleshooting
matched_resources['eventarc'] = matched_resources.get('eventarc', [])
if 'google_eventarc_google_api_source' in missing:
    matched_resources['eventarc'].append({'file': 'GoogleApiSource', 'target': 'google_eventarc_google_api_source'})
if 'google_eventarc_google_channel_config' in missing:
    matched_resources['eventarc'].append({'file': 'GoogleChannelConfig', 'target': 'google_eventarc_google_channel_config'})

matched_resources['dataprocmetastore'] = matched_resources.get('dataprocmetastore', [])
if 'google_dataproc_metastore_service' in missing:
    matched_resources['dataprocmetastore'].append({'file': 'Service', 'target': 'google_dataproc_metastore_service'})

# Clean duplicates
for p in matched_resources:
    unique = {r['target']: r for r in matched_resources[p]}.values()
    matched_resources[p] = list(unique)

with open('matched_by_product.json', 'w') as f:
    json.dump(matched_resources, f, indent=2)

print("Mapped resources written to matched_by_product.json")
