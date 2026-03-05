document.addEventListener("DOMContentLoaded", () => {
  const apiBaseUrl = "http://localhost:8080/api/v1";
  let jwtToken = null;

  const responseArea = document.getElementById("api-response");

  // Helper to display results
  const showResponse = (data) => {
    responseArea.textContent = JSON.stringify(data, null, 2);
  };

  // --- Event Listeners ---

  // Register User
  document
    .getElementById("register-btn")
    .addEventListener("click", async () => {
      const username = document.getElementById("reg-username").value;
      const email = document.getElementById("reg-email").value;
      const password = document.getElementById("reg-password").value;
      const role = document.getElementById("role").value;
      const department = document.getElementById("department").value;

      try {
        const res = await fetch(`${apiBaseUrl}/user/register`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username, email, password, role, department }),
        });
        const data = await res.json();
        showResponse({ status: res.status, body: data });
      } catch (error) {
        showResponse({ error: error.message });
      }
    });

  // Login User
  document.getElementById("login-btn").addEventListener("click", async () => {
    const identifier = document.getElementById("login-identifier").value;
    const password = document.getElementById("login-password").value;

    try {
      const res = await fetch(`${apiBaseUrl}/user/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ identifier, password }),
      });
      const data = await res.json();
      if (data.token) {
        jwtToken = data.token;
        showResponse({
          status: res.status,
          message: "Login successful! Token stored.",
          token: jwtToken,
        });
      } else {
        showResponse({ status: res.status, body: data });
      }
    } catch (error) {
      showResponse({ error: error.message });
    }
  });

  // Create Asset
  document
    .getElementById("create-asset-btn")
    .addEventListener("click", async () => {
      if (!jwtToken) {
        showResponse({ error: "You must be logged in to create an asset." });
        return;
      }

      const name = document.getElementById("asset-name").value;
      const assetType = document.getElementById("asset-type").value;
      const classification = document.getElementById(
        "asset-classification"
      ).value;

      try {
        const res = await fetch(`${apiBaseUrl}/asset/create`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${jwtToken}`,
          },
          body: JSON.stringify({ name, assetType, classification }),
        });
        const data = await res.json();
        showResponse({ status: res.status, body: data });
      } catch (error) {
        showResponse({ error: error.message });
      }
    });

  // Get Asset
  document
    .getElementById("get-asset-btn")
    .addEventListener("click", async () => {
      if (!jwtToken) {
        showResponse({ error: "You must be logged in to get an asset." });
        return;
      }

      const assetId = document.getElementById("get-asset-id").value;
      if (!assetId) {
        showResponse({ error: "Please provide an Asset ID." });
        return;
      }

      try {
        const res = await fetch(`${apiBaseUrl}/asset/${assetId}`, {
          method: "GET",
          headers: {
            Authorization: `Bearer ${jwtToken}`,
          },
        });
        const data = await res.json();
        showResponse({ status: res.status, body: data });
      } catch (error) {
        showResponse({ error: error.message });
      }
    });
  document.querySelectorAll(".scenario-btn").forEach((button) => {
    button.addEventListener("click", async () => {
      if (!jwtToken) {
        showResponse({ error: "Please login first." });
        return;
      }

      const test = button.dataset.test;
      let endpoint = "";
      let options = {};

      switch (test) {
        case "rbac-employee":
          endpoint = `${apiBaseUrl}/asset/create`;
          options = {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${jwtToken}`,
            },
            body: JSON.stringify({
              name: "RBAC Test Asset",
              asset_type: "Test",
              classification: 1,
              owner_username: "emp1",
            }),
          };
          break;
        case "rbac-viewer":
          endpoint = `${apiBaseUrl}/asset/create`;
          options = {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${jwtToken}`,
            },
            body: JSON.stringify({
              name: "RBAC Fail Asset",
              asset_type: "Test",
              classification: 1,
              owner_username: "viewer1",
            }),
          };
          break;
        case "mac-viewer":
          endpoint = `${apiBaseUrl}/asset/doc-ts-001`; // TopSecret doc
          options = {
            method: "GET",
            headers: { Authorization: `Bearer ${jwtToken}` },
          };
          break;
        case "abac-viewer":
          endpoint = `${apiBaseUrl}/asset/server-001`; // IT Server
          options = {
            method: "GET",
            headers: { Authorization: `Bearer ${jwtToken}` },
          };
          break;
        case "dac-viewer":
          endpoint = `${apiBaseUrl}/asset/report-2025-q4`; // Report with direct access
          options = {
            method: "GET",
            headers: { Authorization: `Bearer ${jwtToken}` },
          };
          break;
      }

      try {
        const res = await fetch(endpoint, options);
        const data = await res.json();
        showResponse({ status: res.status, body: data });
      } catch (error) {
        showResponse({ error: error.message });
      }
    });
  });
});
