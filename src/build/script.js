document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("analyze-form");
    const teamLinkInput = document.getElementById("teamLink");
    const loading = document.getElementById("loading");
    const error = document.getElementById("error");
    const report = document.getElementById("report");
    const team = document.getElementById("team");
    const core = document.getElementById("core");
    const mode = document.getElementById("mode");
    const coverage = document.getElementById("coverage");
    const support = document.getElementById("support");

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const teamLink = teamLinkInput.value;

        // Show loading indicator and hide error/report messages
        loading.style.display = "block";
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

            // Populate report fields with the analysis data
            team.textContent = data.team;
            core.textContent = data.core;
            mode.textContent = data.mode;
            coverage.textContent = data.coverage;
            support.textContent = data.support;

            // Show the report and hide the loading indicator
            report.style.display = "block";
        } catch (err) {
            error.textContent = err.message;
            error.style.display = "block";
        } finally {
            loading.style.display = "none";
        }
    });
});
