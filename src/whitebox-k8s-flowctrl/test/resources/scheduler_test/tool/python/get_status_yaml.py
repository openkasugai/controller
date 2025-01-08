#!/usr/bin/env python3
# Copyright 2024 NTT Corporation , FUJITSU LIMITED
import argparse
import sys
import yaml

KEY_STATUS = "status"

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--yamlfile", help="YAML File")
    args = parser.parse_args()

    exitCode = 0
    if args.yamlfile:
        with open(args.yamlfile) as f:
            exitCode = getStatusYaml(f)
    else:
        exitCode = getStatusYaml(sys.stdin)

    return exitCode

def getStatusYaml(f):
    try:
        yml = yaml.safe_load(f)
        if KEY_STATUS in yml:
            yaml.dump({KEY_STATUS: yml.get(KEY_STATUS)}, sys.stdout, indent=2)
            return 0
        elif yml.get("kind") == "List":
            fbasename = f.name.split(".")[0]
            for i, ft in enumerate(yml.get("items")):
                with open(f'{fbasename}_FT{i+1}.yaml', 'w') as file:
                    yaml.dump(ft, file, indent=2)
            return 0
        else:
            print("status element does not exist.")
            return 1
    except Exception as exc:
        print("failed parsing yaml.", exc, file=sys.stderr)
        return False

if __name__ == "__main__":
    sys.exit(main())
