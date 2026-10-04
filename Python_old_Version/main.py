import argparse
from swapi import search_character, get_resource
from utils import format_character, format_homeworld, calculate_time_ratio


def main():
    parser = argparse.ArgumentParser(description="Star Wars CLI Search Tool")
    parser.add_argument("command", choices=["search"])
    parser.add_argument("name", help="Character name")
    parser.add_argument("--world", action="store_true", help="Show homeworld info")

    args = parser.parse_args()

    if args.command == "search":
        results = search_character(args.name)

        if not results:
            print("Character not found.")
            return

        char = results[0]
        print(format_character(char))

        if args.world:
            world = get_resource(char["homeworld"])
            print(format_homeworld(world))
            print(calculate_time_ratio(world))


if __name__ == "__main__":
    main()
