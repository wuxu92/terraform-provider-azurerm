import json
from os import path


def main():
    deprecated = []
    with open(path.join(path.dirname(__file__), './azurerm-schema.json'), 'r') as f:
        schema = json.load(f)
        for rt, resource in schema['providerSchema']['resources'].items():
            if resource.get("deprecaton_message", '') != '':
                deprecated.append(rt)

    print(deprecated)


if __name__ == '__main__':
    main()
