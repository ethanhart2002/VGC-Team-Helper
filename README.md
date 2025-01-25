# VGC-Team-Helper
Website: [VGC Team Helper](https://vgcteamhelper.com/)

Please consider leaving a Watch or a Star if you're interested in following this project!

## Version: 1.1.0
New to this update:
1. Additional links to VGC resources.
2. Fixes UI bug where Paradox Pokemon and Urshifu wouldn't have their images shown in the report.
3. Analysis of your team's offensive capabilities against the top Pokemon in
the current metagame. 
4. Changes to the report to decrease verbosity.
5. Changes to the overall scoring algorithm to reflect the strength of a team against
other Pokemon in the metagame.
6. Changes to the overall scoring algorithm that increases scaling factor in mode scoring, 
core scoring, and support scoring.

## License
See the License [here](./LICENSE).


## Description
This is the open source code for VGC Team Helper, a project that aims to help players build VGC teams with recommendations based on a provided Pokepaste. At a high level, this project does the following:
1. Takes in a Pokepaste link and decodes the entire Pokemon team.
2. Analyzes the Pokemon team for cores, modes, support moves, and coverage moves, all while scoring each aspect. 
3. Returns the score and analysis to the user.


For more information on cores, modes, support moves, and coverage moves, visit [VGC Guide](https://www.vgcguide.com/)

## Installation
1. Ensure Go is installed on your system: [Download Go](https://go.dev/doc/install)
2. Clone this repository to your local machine.
3. Navigate to the project's directory in your clone, and then navigate to the src directory.
4. Visit main.go to change the listening port as directed in func main(), and visit build/script.js to change the hostname from vgcteamhelper.com to the corresponding localhost port as directed in the file.
5. Save the changes and run "go run ." in the src directory.


## Contributing
We welcome contributions. To contribute:
1. Fork this repository off of the main branch. 
2. Create a new branch in your fork. Give it a descriptive name for the change you wish to make. 
3. Follow the steps above in Installation to clone your fork locally and set it up for local testing.
4. Once changes are made, push your branch to your original fork.
5. Submit a pull request into the dev branch of this repo. 


