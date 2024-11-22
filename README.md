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

    Fetches all characters from the SWAPI.
Converts the response into JSON and extracts the results.
2. *def get_char_details(char_url):*

    Fetches detailed properties of a character or a planet from the provided URL.

3. *def search_char_name(results, name):*

    Searches for a specific character by name from the list of fetched characters.
4. *Main*  

    Uses argparse to handle CLI commands and flags.
Displays character and homeworld details based on user input.
