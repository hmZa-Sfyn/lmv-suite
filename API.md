# LanManVan Module Manager

A fast, lightweight FastAPI-based module manager for LanManVan security & utility tools.  
Stores metadata in SQLite, handles git operations, tracks usage stats, and provides a clean JSON API.

**Live URL**: http://194.164.194.118:10382  
**Interactive Docs**: http://194.164.194.118:10382/docs (Swagger UI)  
**Redoc Docs**: http://194.164.194.118:10382/redoc

---

### Description

Personal registry for 20 pre-seeded LanManVan modules (port scanner, hash cracker, subdomain enum, etc.).  
Allows browsing, adding custom modules, downloading via git, tracking views/downloads, and user account management.

Publicly accessible on port 10382 (listening on 0.0.0.0).

---

### Use Cases

- Browse and download your tool collection from any device.
- Keep tools updated with one command (`git pull` automatically).
- Add your own private/custom git-based tools.
- Track which modules you've already cloned on the server.
- Integrate into CLI scripts, web dashboards, or automation workflows.

---

### Authentication

All endpoints except `/auth/*` require a JWT token.

<code>
# Register (once)
curl -X POST http://194.164.194.118:10382/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "myuser", "password": "StrongPass123!"}'
</code>

<code>
# Login → get token
curl -X POST http://194.164.194.118:10382/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "myuser", "password": "StrongPass123!"}'
</code>

Use token in header for all other calls:

<code>
Authorization: Bearer YOUR_TOKEN_HERE
</code>

---

### API Endpoints

#### **POST /auth/register**
- **Description**: Create a new user account
- **Body Params**:
  - username (string, required)
  - password (string, required)
- **Example**:
<code>
curl -X POST http://194.164.194.118:10382/auth/register \
  -d '{"username": "alice", "password": "Secret123!"}'
</code>

#### **POST /auth/login**
- **Description**: Login and receive JWT token (valid 7 days)
- **Body Params**:
  - username (string, required)
  - password (string, required)
- **Example**:
<code>
curl -X POST http://194.164.194.118:10382/auth/login \
  -d '{"username": "alice", "password": "Secret123!"}'
</code>

#### **GET /modules**
- **Description**: List all available modules
- **Auth**: Required
- **Example**:
<code>
curl -H "Authorization: Bearer $TOKEN" http://194.164.194.118:10382/modules
</code>

#### **GET /modules/{module_id}**
- **Description**: Get single module details (increments view_count)
- **Path Param**: module_id (int)
- **Auth**: Required
- **Example**:
<code>
curl -H "Authorization: Bearer $TOKEN" http://194.164.194.118:10382/modules/5
</code>

#### **POST /modules**
- **Description**: Add a new module
- **Body Params**:
  - slug (string, required, unique)
  - name (string, required)
  - description (string, required)
  - repo_url (string, optional)
  - version (string, optional, default "latest")
- **Auth**: Required
- **Example**:
<code>
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -d '{"slug": "my-tool", "name": "My Tool", "description": "Custom tool", "repo_url": "https://github.com/user/tool.git"}' \
  http://194.164.194.118:10382/modules
</code>

#### **PATCH /modules/{module_id}**
- **Description**: Partially update a module (e.g., change repo_url or version)
- **Path Param**: module_id (int)
- **Body Params**: any of name, description, repo_url, version
- **Auth**: Required
- **Example**:
<code>
curl -X PATCH -H "Authorization: Bearer $TOKEN" \
  -d '{"repo_url": "https://new-repo.com/tool.git"}' \
  http://194.164.194.118:10382/modules/10
</code>

#### **DELETE /modules/{module_id}**
- **Description**: Delete module and its local repo folder
- **Path Param**: module_id (int)
- **Auth**: Required
- **Example**:
<code>
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://194.164.194.118:10382/modules/15
</code>

#### **POST /modules/{module_id}/download**
- **Description**: Git clone (first time) or pull (updates) the module
- **Path Param**: module_id (int)
- **Auth**: Required
- **Example**:
<code>
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://194.164.194.118:10382/modules/3/download
</code>

#### **GET /modules/{module_id}/repo**
- **Description**: Return only the repo_url (increments download_count)
- **Path Param**: module_id (int)
- **Auth**: Required
- **Example**:
<code>
curl -H "Authorization: Bearer $TOKEN" \
  http://194.164.194.118:10382/modules/7/repo
</code>

#### **GET /my/repos**
- **Description**: List only modules you have downloaded/cloned
- **Auth**: Required
- **Example**:
<code>
curl -H "Authorization: Bearer $TOKEN" http://194.164.194.118:10382/my/repos
</code>

#### **GET /settings**
- **Description**: View your current username
- **Auth**: Required
- **Example**:
<code>
curl -H "Authorization: Bearer $TOKEN" http://194.164.194.118:10382/settings
</code>

#### **PATCH /settings**
- **Description**: Change username and/or password
- **Body Params**: username (optional), password (optional)
- **Auth**: Required
- **Example**:
<code>
curl -X PATCH -H "Authorization: Bearer $TOKEN" \
  -d '{"password": "NewPass2025!"}' \
  http://194.164.194.118:10382/settings
</code>

---

### Auto CLI Script

Save as <code>lmv.sh</code> (edit USER/PASS):

<code>
#!/bin/bash
API="http://194.164.194.118:10382"
USER="myuser"
PASS="StrongPass123!"

response=$(curl -s -X POST "$API/auth/login" -H "Content-Type: application/json" -d "{\"username\":\"$USER\",\"password\":\"$PASS\"}")
TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

[ -z "$TOKEN" ] && echo "Login failed!" && exit 1

echo "Logged in!"
echo -e "\nAll Modules:"
curl -s -H "Authorization: Bearer $TOKEN" "$API/modules" | jq -r '.[] | "\(.id) | \(.slug | ljust(20)) | \(.name)"' 2>/dev/null || echo "Install jq for pretty output"

echo -e "\nDownloaded Repos:"
curl -s -H "Authorization: Bearer $TOKEN" "$API/my/repos" | jq
</code>

Run:
<code>
chmod +x lmv.sh && ./lmv.sh
</code>

---

Enjoy your global LanManVan module hub! 🚀
