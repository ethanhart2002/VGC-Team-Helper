
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



    // Function to replace newlines with <br> tags
    function handleNewlines(text) {
        return text.replace(/\n/g, '<li>');
    }

    function reportCoverage(coverageArray) {

        if (coverageArray.length <= 0) {
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


    const nameDict = {
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
        //TODO Add Galarian, Alolan, and Paldean conversions
    };

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const teamLink = teamLinkInput.value;

        // Show loading indicator and hide error/report messages
        document.getElementById("loading").style.display = "block";
        error.style.display = "none";
        report.style.display = "none";

        try {
            const response = await fetch("http://localhost:8080/analyze", {
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
