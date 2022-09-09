import os
import shutil
import sys
import yaml


def main():
  file_name = sys.argv[1]
  shutil.copy(file_name, os.getcwd())
  samples_dir = os.path.split(file_name)[0]
  sample = yaml.safe_load(open(file_name))
  print(yaml.dump(sample))
  all_deps = [sample['resource']]
  updates = sample.get('updates')
  if updates:
    all_deps.append(updates[0]['resource'])
  all_deps += sample.get('dependencies', [])
  for dep in all_deps:
    dep_file_name = f'{samples_dir}/{os.path.split(dep)[1]}'
    shutil.copy(dep_file_name, os.getcwd())


if __name__ == '__main__':
  main()