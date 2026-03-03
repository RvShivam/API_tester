/* ============================================================
   app.js — API Tester Web Terminal Client
   Uses xterm.js for terminal rendering and WebSocket for
   relaying commands to the Go backend.
   ============================================================ */

(function () {
    "use strict";

    // ── 1. Resolve WebSocket URL automatically ──────────────────
    // In production, Render serves HTTPS so we use wss://.
    // Locally, we use ws://.
    const proto = location.protocol === "https:" ? "wss" : "ws";
    const WS_URL = `${proto}://${location.host}/ws`;

    // ── 2. Initialise xterm.js ──────────────────────────────────
    const term = new Terminal({
        cursorBlink: true,
        fontFamily: "'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace",
        fontSize: 13,
        lineHeight: 1.5,
        theme: {
            background: "#010409",
            foreground: "#e6edf3",
            cursor: "#58a6ff",
            cursorAccent: "#0d1117",
            selectionBackground: "rgba(88,166,255,0.25)",
            black: "#484f58", red: "#ff7b72",
            green: "#3fb950", yellow: "#d29922",
            blue: "#58a6ff", magenta: "#bc8cff",
            cyan: "#76e3ea", white: "#b1bac4",
            brightBlack: "#6e7681", brightRed: "#ffa198",
            brightGreen: "#56d364", brightYellow: "#e3b341",
            brightBlue: "#79c0ff", brightMagenta: "#d2a8ff",
            brightCyan: "#b3f0ff", brightWhite: "#f0f6fc",
        },
        scrollback: 2000,
        allowProposedApi: true,
    });

    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    term.open(document.getElementById("terminal"));
    fitAddon.fit();

    window.addEventListener("resize", () => fitAddon.fit());

    // ── 3. WebSocket setup ──────────────────────────────────────
    let ws = null;
    let reconnectTimer = null;
    const connStatusEl = document.getElementById("conn-status");
    const runBtn = document.getElementById("run-btn");

    function setStatus(state) {
        connStatusEl.className = "";
        if (state === "connected") {
            connStatusEl.className = "conn-connected";
            connStatusEl.textContent = "● Connected";
            runBtn.disabled = false;
        } else if (state === "error") {
            connStatusEl.className = "conn-error";
            connStatusEl.textContent = "✖ Disconnected — retrying…";
            runBtn.disabled = true;
        } else {
            connStatusEl.className = "conn-connecting";
            connStatusEl.textContent = "⟳ Connecting…";
            runBtn.disabled = true;
        }
    }

    function connect() {
        setStatus("connecting");
        ws = new WebSocket(WS_URL);

        ws.onopen = () => {
            setStatus("connected");
            clearTimeout(reconnectTimer);
        };

        ws.onmessage = (event) => {
            // xterm.js uses \r\n for line endings — the server already handles this.
            term.write(event.data);
        };

        ws.onclose = () => {
            setStatus("error");
            reconnectTimer = setTimeout(connect, 3000);
        };

        ws.onerror = () => {
            ws.close();
        };
    }

    connect();

    // ── 4. Command execution ─────────────────────────────────────
    const cmdInput = document.getElementById("cmd-input");
    let lastCommand = "";

    function runCommand(cmd) {
        cmd = cmd.trim();
        if (!cmd || !ws || ws.readyState !== WebSocket.OPEN) return;
        lastCommand = cmd;
        cmdInput.value = "";
        ws.send(cmd);
    }

    // Run button
    runBtn.addEventListener("click", () => runCommand(cmdInput.value));

    // Enter key to submit
    cmdInput.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            runCommand(cmdInput.value);
        } else if (e.key === "ArrowUp") {
            // Recall last command (like a real terminal)
            if (lastCommand) {
                cmdInput.value = lastCommand;
                // Move cursor to end
                setTimeout(() => cmdInput.setSelectionRange(cmdInput.value.length, cmdInput.value.length), 0);
            }
        }
    });

    // ── 5. Demo command buttons ──────────────────────────────────
    document.querySelectorAll(".demo-btn").forEach((btn) => {
        btn.addEventListener("click", () => {
            // Support both simple `data-command` and structured `data-cmd` + optional flags.
            // Using separate attributes for body/headers/auth avoids HTML quoting hell with JSON.
            let cmd = btn.dataset.command || "";

            if (!cmd && btn.dataset.cmd) {
                cmd = btn.dataset.cmd;
                if (btn.dataset.body) cmd += ` --body '${btn.dataset.body}'`;
                if (btn.dataset.headers) cmd += ` --headers "${btn.dataset.headers}"`;
                if (btn.dataset.auth) cmd += ` --auth "${btn.dataset.auth}"`;
                if (btn.dataset.extra) cmd += ` ${btn.dataset.extra}`;
            }

            if (!cmd) return;

            // Visual feedback: briefly mark the button as running.
            btn.classList.add("running");
            setTimeout(() => btn.classList.remove("running"), 1500);

            // Populate the text box so the user can see what ran.
            cmdInput.value = cmd;
            cmdInput.focus();

            // Execute immediately.
            runCommand(cmd);
        });
    });

    // ── 6. Clear button ──────────────────────────────────────────
    document.getElementById("clear-btn").addEventListener("click", () => {
        term.clear();
    });

    // ── 7. Focus input on page click ────────────────────────────
    document.querySelector(".terminal-area").addEventListener("click", (e) => {
        if (e.target.tagName !== "BUTTON") {
            cmdInput.focus();
        }
    });

})();
