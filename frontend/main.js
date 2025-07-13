
// Store competitions by ID for updates
let competitionsState = {};

function renderCompetitions(dataArr, highlightUser = null, compId = null, highlightScore = null) {
    const container = document.getElementById("competitions");
    container.innerHTML = "";
    dataArr.forEach(event => {
        // Sort users by score descending
        let table = `<div class='competition-table'><h2>Competition ${event.CompetitionID}</h2><table><thead><tr><th style='text-align:center; width:90px;'>User</th><th style='text-align:center;'>Score</th></tr></thead><tbody>`;
        event.Users.forEach(user => {
            let highlightUserClass = (highlightUser && compId === event.CompetitionID && user.id === highlightUser) ? "highlight-row" : "";
            let scoreFormatted = `$ ${parseFloat(user.score).toFixed(2)}`;
            let cellId = `score-${event.CompetitionID}-${user.id}`;
            table += `<tr class='${highlightUserClass}'><td style='text-align:center; width:90px;'>${user.id}</td><td style='text-align:center;' id='${cellId}'>${scoreFormatted}</td></tr>`;
        });
        table += "</tbody></table></div>";
        container.innerHTML += table;
    });
    // After rendering, add highlight class to all updated cells
    if (Array.isArray(highlightScore) && compId) {
        highlightScore.forEach(userId => {
            const cell = document.getElementById(`score-${compId}-${userId}`);
            if (cell) {
                cell.classList.remove("highlight-cell"); // Remove if present
                // Force reflow to restart animation
                void cell.offsetWidth;
                cell.classList.add("highlight-cell");
            }
        });
    } else if (highlightScore && compId) {
        const cell = document.getElementById(`score-${compId}-${highlightScore}`);
        if (cell) {
            cell.classList.remove("highlight-cell");
            void cell.offsetWidth;
            cell.classList.add("highlight-cell");
        }
    }
}

// Initial render with dummy data
renderCompetitions(Object.values(competitionsState));

// WebSocket logic
let ws;
function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = function() {
        console.log("Connected to WebSocket server");
    };

    ws.onmessage = function(event) {
        try {
            const msg = JSON.parse(event.data);
            // msg is expected to be {CompetitionID, Users}
            const compId = msg.CompetitionID;
            let highlightUser = null;
            let highlightScore = [];
            if (competitionsState[compId]) {
                // Find updated user(s)
                const oldUsers = competitionsState[compId].Users;
                msg.Users.forEach(newUser => {
                    const oldUser = oldUsers.find(u => u.id === newUser.id);
                    if (!oldUser || oldUser.score !== newUser.score) {
                        highlightUser = newUser.id;
                        highlightScore.push(newUser.id);
                    }
                });
            }
            competitionsState[compId] = msg;
            renderCompetitions(Object.values(competitionsState), highlightUser, compId, highlightScore.length ? highlightScore : null);
        } catch (e) {
            console.error("Invalid message format", e);
        }
    };

    ws.onclose = function() {
        console.log("WebSocket connection closed, retrying...");
        setTimeout(connect, 1000); // Reconnect after 1 second
    };

    ws.onerror = function(error) {
        console.error("WebSocket error:", error);
    };
}

connect();
