# SecureAsset – Secure Asset Management System

SecureAsset is a backend-only, security-focused asset management system developed using **Go (Gin)** and **MongoDB**.  
The system provides APIs for securely managing organizational assets while enforcing strict authentication, authorization, and auditing policies.

This project is designed for academic and practical demonstration of **comprehensive access-control mechanisms**, including:

- Mandatory Access Control (MAC)  
- Discretionary Access Control (DAC)  
- Role-Based Access Control (RBAC)  
- Rule-Based Access Control (RuBAC)  
- Attribute-Based Access Control (ABAC)

All development and testing are performed through **RESTful APIs** using **Postman**.

---

## 🚀 Project Overview

SecureAsset allows organizations to:

- Register and authenticate users  
- Manage assets (equipment, devices, tools, digital assets)  
- Request, approve, borrow, return, and track assets  
- Enforce strict access policies based on roles, attributes, rules, and classifications  
- Log all activities and system operations for security review  
- Protect data with encrypted logs, secure password storage, and MFA  

---

## 🔐 Key Security Features

### **1. Access Control Models**
#### 🔸 Mandatory Access Control (MAC)
- Assets classified as *Public*, *Internal*, *Restricted*, *Confidential*  
- System admins assign clearance levels  
- Users cannot override security labels  

#### 🔸 Discretionary Access Control (DAC)
- Asset owners (department heads / managers) can grant or revoke access  
- Permission logs track all DAC changes  

#### 🔸 Role-Based Access Control (RBAC)
Roles include:
- Admin  
- Asset Manager  
- Department Head  
- Employee  

Each role has fixed permissions enforced automatically by the system.

#### 🔸 Rule-Based Access Control (RuBAC)
Examples:
- Deny asset requests outside working hours  
- Restrict high-value asset access to office network  
- Allow only IT managers to approve server-related assets  

#### 🔸 Attribute-Based Access Control (ABAC)
Attributes used:
- Role  
- Department  
- Location  
- Clearance level  
- Time of access  
- Asset sensitivity  

Dynamic attribute evaluation determines final access decisions.

---

## 🔑 Identification & Authentication

- Secure user registration  
- Email/phone verification  
- CAPTCHA for bot protection  
- Strong password policy enforcement  
- Password hashing (bcrypt/argon2)  
- JWT / token-based authentication  
- MFA using OTP  
- Account lockout after repeated failed logins  
- Secure session lifecycle management  

---

## 📝 Audit Trails & Logging

- Logs every user action (create, update, view, delete)  
- Logs system events (startup, shutdown, config changes)  
- Logs include:  
  - Username  
  - Action  
  - Endpoint  
  - Timestamp  
  - IP address  
- Logs stored encrypted  
- Centralized log aggregation  
- Alerts for:
  - Failed login attempts  
  - Access to restricted assets  
  - Permission modification events  

---

## 💾 Database & Backups

- MongoDB as primary database  
- Encrypted database fields for sensitive data  
- Automated periodic backups  
- Backup recovery strategy included  

---

## 🧱 Tech Stack

- **Language:** Go (Golang)  
- **Framework:** Gin Web Framework  
- **Database:** MongoDB  
- **Auth:** JWT, OTP (MFA)  
- **Testing:** Postman  


