# Access Control Implementation Guide

Implementing multiple access control models can be complex. The `github.com/casbin/casbin/v2` library is highly recommended as it is a powerful and flexible authorization library that can handle all the required models (MAC, DAC, RBAC, RuBAC, ABAC) with a single, unified API.

You define your access control model in a `.conf` file and your policies in a `.csv` file or a database.

## 1. The Casbin Model (`model.conf`)

Here is a comprehensive model that combines RBAC with ABAC, which can be extended to cover all your requirements. Create a file named `model.conf` in your `internal/config` directory.

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act || r.sub == "admin"
```

### Explanation:

- **`[request_definition]`**: Defines the structure of an access request. Here, it's `subject` (user), `object` (resource), and `action` (e.g., read, write).
- **`[policy_definition]`**: Defines the structure of a policy rule. `sub` can be a user or a role. `eft` (effect) allows for both `allow` and `deny` rules.
- **`[role_definition]`**: Defines the role hierarchy (RBAC). `g = _, _` means a user can be in a role.
- **`[policy_effect]`**: Determines the final decision. It allows access if at least one policy allows it and no policy denies it.
- **`[matchers]`**: The core logic.
  - `g(r.sub, p.sub)`: Checks if the user (`r.sub`) has the required role (`p.sub`).
  - `keyMatch(r.obj, p.obj)`: Matches the requested resource with policy resources (supports wildcards like `/assets/*`).
  - `r.act == p.act`: Matches the requested action.
  - `r.sub == "admin"`: A simple rule to grant the "admin" user universal access.

## 2. Implementing the Access Control Models

### 🔸 Role-Based Access Control (RBAC)

RBAC is directly supported by the `[role_definition]` section.

**Policy Example (`policy.csv`):**

```csv
p, asset_manager, /api/v1/assets/*, (GET|POST)
p, employee, /api/v1/assets/me, GET

g, alice, asset_manager
g, bob, employee
```

- The first line gives `asset_manager` roles GET and POST access to all asset routes.
- The second line gives `employee` roles GET access to their own assets.
- `g` lines assign users `alice` and `bob` to their respective roles.

### 🔸 Attribute-Based (ABAC) & Rule-Based (RuBAC)

For ABAC/RuBAC, you can pass custom objects into the matcher. The matcher function can then evaluate attributes from these objects.

**Extended `model.conf` Matcher for ABAC/RuBAC:**

```ini
[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act && timeCheck(r.sub.Time) && locationCheck(r.sub.Location)
```

Here, `r.sub` would be a struct containing user attributes like `Time` and `Location`. You would need to implement the `timeCheck` and `locationCheck` functions and add them to the Casbin enforcer.

### 🔸 Discretionary Access Control (DAC)

DAC allows resource owners to grant permissions. You can implement this by adding API endpoints that allow authorized users (owners) to add or remove Casbin policies.

**Example:** An asset owner could make a `POST /api/v1/assets/{id}/permissions` request, and your handler would call `casbinEnforcer.AddPolicy(...)` to grant another user access.

### 🔸 Mandatory Access Control (MAC)

MAC is about data classification and user clearance. You can model this by adding classification levels to your objects and clearance levels to your subjects.

**Policy Example for MAC:**

```csv
# p, subject_clearance, object_classification, action
p, confidential, confidential, (read|write)
p, restricted, restricted, (read|write)
p, internal, internal, (read|write)
```

Your application logic would then check if a user's clearance level is sufficient for the asset's classification level before checking the standard Casbin policy.

## 3. Storing Policies

For production, you should store policies in your MongoDB database using the `github.com/casbin/mongo-adapter/v3`. This allows you to manage policies dynamically without restarting the application.
