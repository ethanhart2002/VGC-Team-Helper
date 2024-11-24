
document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("analyze-form");
    const teamLinkInput = document.getElementById("teamLink");
    const loading = document.getElementById("loading");
    const error = document.getElementById("error");
    const report = document.getElementById("report");
    const core = document.getElementById("core");
    const mode = document.getElementById("mode");
    const coverage = document.getElementById("coverage");
    const support = document.getElementById("support");
    const score = document.getElementById("score");
    const tablinks = document.querySelectorAll(".tablinks");
    const tabcontents = document.querySelectorAll(".tabcontent");


    tablinks.forEach((button) => {
        button.addEventListener("click", function () {
            const target = document.querySelector(`#${this.getAttribute("data-tab")}`);

            // Toggle the active class on the clicked button and its content
            const isActive = this.classList.contains("active");

            // Hide all other panels and remove active state from other buttons
            tablinks.forEach((btn) => btn.classList.remove("active"));
            tabcontents.forEach((content) => content.classList.remove("active"));

            // If the clicked button wasn't already active, activate it
            if (!isActive) {
                this.classList.add("active");
                if (target) {
                    target.classList.add("active");
                }
            }
        });
    });



    // Function to replace newlines with <br> tags
    function handleNewlines(text) {
        return text.replace(/\n/g, '<li>');
    }

    function reportCoverage(coverageArray) {

        if (coverageArray == null) {
            coverage.innerHTML = `
            <li> Your team has coverage options to hit all 18 types! </li>
            `
        } else {
            let rep = `<li> Your team is missing attacking moves that can hit the following types for super-effective damage: </li>`;
            for (let i = 0; i < coverageArray.length; i++) {
                rep += `<br><b> ${coverageArray[i]} </b>`;
            }
            rep += `<br><br><li> If you find your team is struggling against Pokemon of these types, considering adding coverage moves to hit these Pokemon with super effective damage. </li>`;
            coverage.innerHTML = rep;
        }
    }

    function reportScore(scoreInput) {
        score.innerHTML = ``;
        score.className = "";
        let s = ``
        if (typeof scoreInput === 'number' && !isNaN(scoreInput)) {
            if (scoreInput >= 8.0) {
                score.classList.add("good");
            } else if (scoreInput >= 7.0) {
                score.classList.add("okay");
            } else {
                score.classList.add("needsWork");
            }
            s += `${scoreInput.toString()}/10`;
        } else {
            s += `Score not available`;
        }
        score.innerHTML = s;
    }


    function parseStructs(teamMembers) {
        const container = document.getElementById("team-container");

        // Clear previous content
        container.innerHTML = '';

        teamMembers.forEach(struct => {
            const structElement = document.createElement("div");
            structElement.classList.add("pokemon-container"); // Add class for styling

            const pokemonName = struct.pokemon in nameDict ? nameDict[struct.pokemon] : struct.pokemon;
            const spriteUrl = `https://img.pokemondb.net/sprites/scarlet-violet/normal/${pokemonName.toLowerCase()}.png`;

            const img = document.createElement('img');
            img.src = spriteUrl;
            img.alt = struct.pokemon;
            img.onerror = () => { img.style.display = 'none'; }; // Hide image if not found
            img.classList.add("pokemon-image")

            const textContainer = document.createElement('div');
            textContainer.classList.add("pokemon-textbox"); // Add class for styling
            textContainer.innerHTML = `
            <p>Name: ${struct.pokemon}</p>
            <p>Type(s): ${struct.type}</p>
            <p>Item: ${struct.item}</p>
            <p>Ability: ${struct.ability}</p>
            <p>Tera Type: ${struct.tera_type}</p>
            <p>Moves: ${struct.moves}</p>
        `;

            structElement.appendChild(img);
            structElement.appendChild(textContainer);
            container.appendChild(structElement);
        });
    }



    // There is a discrepancy between how Pokepaste names certain Pokemon vs PokemonDB, so
    // we will need to convert them sometimes.
    const nameDict = {
        "Pyroar-M": "Pyroar-Male",
        "Pyroar-F": "Pyroar-Female",
        "Indeedee-F": "Indeedee-female",
        "Basculegion-F": "Basculegion-female",
        "Arcanine-Hisui": "Arcanine-Hisuian",
        "Growlithe-Hisui": "Growlithe-Hisuian",
        "Voltorb-Hisui": "Voltorb-Hisuian",
        "Electrode-Hisui": "Electrode-Hisuian",
        "Typhlosion-Hisui": "Typhlosion-Hisuian",
        "Qwilfish-Hisui": "Qwilfish-Hisuian",
        "Sneasel-Hisui": "Sneasel-Hisuian",
        "Lilligant-Hisui": "Lilligant-Hisuian",
        "Zorua-Hisui": "Zorua-Hisuian",
        "Zoroark-Hisui": "Zoroark-Hisuian",
        "Samurott-Hisui": "Samurott-Hisuian",
        "Braviary-Hisui": "Braviary-Hisuian",
        "Sliggoo-Hisui": "Sliggoo-Hisuian",
        "Goodra-Hisui": "Goodra-Hisuian",
        "Avalugg-Hisui": "Avalugg-Hisuian",
        "Wooper-Paldea": "Wooper-Paldean",
        "Tauros-Paldea-Aqua": "Tauros-Paldean-Aqua",
        "Tauros-Paldea-Blaze": "Tauros-Paldean-Blaze",
        "Tauros-Paldea-Combat": "Tauros-Paldean-Combat",
        "Rattata-Alola": "Rattata-Alolan",
        "Raticate-Alola": "Raticate-Alolan",
        "Raichu-Alola": "Raichu-Alolan",
        "Sandshrew-Alola": "Sandshrew-Alolan",
        "Sandslash-Alola": "Sandslash-Alolan",
        "Vulpix-Alola": "Vulpix-Alolan",
        "Ninetales-Alola": "Ninetales-Alolan",
        "Diglett-Alola": "Diglett-Alolan",
        "Dugtrio-Alola": "Dugtrio-Alolan",
        "Meowth-Alola": "Meowth-Alolan",
        "Persian-Alola": "Persian-Alolan",
        "Geodude-Alola": "Geodude-Alolan",
        "Graveler-Alola": "Graveler-Alolan",
        "Golem-Alola": "Golem-Alolan",
        "Grimer-Alola": "Grimer-Alolan",
        "Muk-Alola": "Muk-Alolan",
        "Exeggutor-Alola": "Exeggutor-Alolan",
        "Marowak-Alola": "Marowak-Alolan",
        "Meowth-Galar": "Meowth-Galarian",
        "Ponyta-Galar": "Ponyta-Galarian",
        "Rapidash-Galar": "Rapidash-Galarian",
        "Slowpoke-Galar": "Slowpoke-Galarian",
        "Slowbro-Galar": "Slowbro-Galarian",
        "Weezing-Galar": "Weezing-Galarian",
        "Farfetch'd-Galar": "Farfetch'd-Galarian",
        "Mr.Mime-Galar": "Mr.Mime-Galarian",
        "Articuno-Galar": "Articuno-Galarian",
        "Zapdos-Galar": "Zapdos-Galarian",
        "Moltres-Galar": "Moltres-Galarian",
        "Slowking-Galar": "Slowking-Galarian",
        "Corsola-Galar": "Corsola-Galarian",
        "Zigzagoon-Galar": "Zigzagoon-Galarian",
        "Linoone-Galar": "Linoone-Galarian",
        "Darumaka-Galar": "Darumaka-Galarian",
        "Darmanitan-Galar": "Darmanitan-Galarian",
        "Yanmask-Galar": "Yanmask-Galarian",
        "Stunfisk-Galar": "Stunfisk-Galarian",
        //TODO Add more naming convention fixes when they come up
    };




    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const teamLink = teamLinkInput.value;

        // Show loading indicator and hide error/report messages
        document.getElementById("loading").style.display = "block";
        error.style.display = "none";
        report.style.display = "none";

        let hostPath = "https://vgcteamhelper.com/analyze";

        //local debug
        //hostPath = "http://localhost:8080/analyze";

        try {
            const response = await fetch(hostPath, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ link: teamLink })
            });

            if (!response.ok) {
                throw new Error("Failed to analyze team. Please check the link.");
            }

            const data = await response.json();


            parseStructs(data.team);
            core.innerHTML = handleNewlines(data.core);
            mode.innerHTML = handleNewlines(data.mode);
            reportCoverage(data.coverage)
            support.innerHTML = handleNewlines(data.support);
            reportScore(data.score);


            // Show the report and hide the loading indicator
            report.style.display = "block";
        } catch (err) {
            error.textContent = err.message;
            error.style.display = "block";
        } finally {
            document.getElementById("loading").style.display = "none";
        }
    });
});
