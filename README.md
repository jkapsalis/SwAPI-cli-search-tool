# Star Wars API cli search tool

This is a cli tool to search for Star Wars characters to find 
information about their height, mass and also about their homeworld details 
using the SWAPI (star wars API).

----------------------------------
#### Features
1. *Character Search*
* Gets and print information about a Star Wars character the user searched.
* It has attributes such as height, mass, and birth year of the character.
2. *Homeworld Details*
* Prints the character's homeworld name, population, and orbital - rotation.
* Calculates the time on the homeworld's character relates to Earth years and days and prints an ratio.

--------------------------------
#### Prerequisites

Install the __requirements.txt__ at my git repository :)

----------------------
#### How to Run
Clone the repository and follow the commands to run the script   :)

    ###for the character search###
    python main.py search 'luke sky'

    ###for the homeworld details search of that character###
    python main.py search'luke sky' --world
--------------------------------------------
#### Sample Outputs

![image](https://github.com/user-attachments/assets/fad76912-52ec-4ec6-b697-c379a46cd911)


![image](https://github.com/user-attachments/assets/fbecc80e-d2d5-469a-8b53-91c287b86033)

---------------------------------------
### Code Overview

1. *def get_all_chars():*

    Gets all the characters from the SWAPI.
Converts the data response into JSON and extract all the info.
2. *def get_char_details(char_url):*

    Gets the properties of a character or a planet from the URL.

3. *def search_char_name(results, name):*

    Searches for a specific character by name from the list of fetched characters.
4. *Main*  

    Uses argparse to handle CLI commands and flags.
Prints character and homeworld information based on the user's input.



-------------------------------------------------------

## Next Golang Objectives (Future Improvements)

As part of expanding this project and improving my backend development skills, I plan to reimplement and enhance this CLI tool using Golang. 

 The objectives are:

### 1. Build a CLI Tool in Go
- Recreate the current Python CLI functionality using Go.
- Use packages like `flag` or `cobra` for command-line argument parsing.
- Maintain similar commands:
  - Character search
  - Homeworld details (`--world` flag equivalent)

### 2. Work with HTTP Requests
- Use Go’s `net/http` package to fetch data from the SWAPI.
- Handle API responses efficiently and explore concurrency where useful.

### 3. JSON Parsing and Structs
- Define Go structs for:
  - Characters
  - Planets (homeworld data)
- Practice unmarshalling JSON responses into typed structs.

### 4. Improve Performance
- Optimize API calls (e.g., caching results or reducing redundant requests).
- Use goroutines to fetch character and homeworld data in parallel.

### 5. Error Handling
- Implement robust error handling using Go’s explicit error system.
- Handle edge cases such as:
  - Character not found
  - API errors or downtime
  - Invalid input

### 6. Modular Code Design
- Organize the project into packages (e.g., `api`, `models`, `cli`).
- Improve code readability and maintainability.

### 7. Cross-Platform CLI Distribution
- Compile the Go program into a standalone binary.
- Allow users to run the tool without installing dependencies.

### 8. Testing
- Write unit tests for:
  - API calls
  - Search functionality
- Use Go’s built-in `testing` package.
