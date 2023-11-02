# Conjur in GoLang

This is an effort to learn enough GoLang to be dangerous and explore Casbin as an RBAC engine.



## Architecture

```plantuml
left to right direction

actor "Client" as client

component "Identity" as identity

component "Authentication" as authentication {
  component "User Login" as user_login
  component "User/Host API Key" as host_user_login
  component "Azure Authenticator" as azure
  component "Google Cloud Authenticator" as gcp
  component "AWS Authenticator" as aws
  component "LDAP Authenticator" as ldap
  component "OIDC Authenticator" as oidc
  component "JWT Authenticator" as jwt
}

component "Authorization" as authorization
component "Policy Loader" as policy_loader
component "Secret Store" as secret_store

client --> authentication : (1) Establishes identity using
client --> authorization : (2) Uses identity to request secret using
client --> secret_store : (3) Recieves secret if authorized

client --> policy_loader : applies changes to
policy_loader --> authentication : Defines Authenticators and\nallows roles to use
policy_loader --> secret_store : Defines variables stored in
policy_loader --> authorization : Defines role permissions for
```

Current Schema:
```
 Schema |         Name          | Type  | Owner       Purp
--------+-----------------------+-------+--------     ----------
 public | annotations           | table | conjur      Identity
 public | authenticator_configs | table | conjur      Authentication
 public | credentials           | table | conjur      Authentication
 public | host_factory_tokens   | table | conjur      Identity (questionable...)
 public | permissions           | table | conjur      Authorization
 public | policy_log            | table | conjur      Policy Loader
 public | policy_versions       | table | conjur      Policy Loader
 public | resources             | table | conjur      Secret Store
 public | resources_textsearch  | table | conjur      Secret Store
 public | role_memberships      | table | conjur      Authorization
 public | roles                 | table | conjur      Authorization
 public | schema_info           | table | conjur      System
 public | schema_migrations     | table | conjur      System
 public | secrets               | table | conjur      Secret Store
 public | slosilo_keystore      | table | conjur      Authentication
```


## Authorization

```
conjur=# select * from role_memberships limit 10;
                          role_id                          |                 member_id                  | admin_option | ownership |               policy_id
-----------------------------------------------------------+--------------------------------------------+--------------+-----------+---------------------------------------
 demo:user:admin                                           | demo:user:admin                            | t            | t         |
 !:!:root                                                  | demo:user:admin                            | t            | f         |
 system:user:admin                                         | demo:user:admin                            | t            | t         |
 system:policy:root                                        | system:user:admin                          | t            | t         |
 system:user:admin                                         | demo:user:admin                            | f            | f         | system:policy:root
 system:policy:conjur                                      | system:user:admin                          | t            | t         | system:policy:root
 system:policy:conjur/replication-sets                     | system:policy:conjur                       | t            | t         | system:policy:conjur
 system:policy:conjur/replication-sets/full                | system:policy:conjur/replication-sets      | t            | t         | system:policy:conjur/replication-sets
 system:group:conjur/replication-sets/full/replicated-data | system:policy:conjur/replication-sets/full | t            | t         | system:policy:conjur/replication-sets
 demo:policy:root                                          | demo:user:admin                            | t            | t         |
(10 rows)

conjur=# select * from roles limit 10;
                          role_id                          |         created_at         |               policy_id
-----------------------------------------------------------+----------------------------+---------------------------------------
 !:!:root                                                  | 2023-06-12 19:32:32.318802 |
 demo:user:admin                                           | 2023-06-12 19:32:32.332894 |
 system:user:admin                                         | 2023-06-12 19:33:04.709486 |
 system:policy:root                                        | 2023-06-12 19:33:04.809451 |
 system:policy:conjur                                      | 2023-06-12 19:33:04.996902 | system:policy:root
 system:policy:conjur/replication-sets                     | 2023-06-12 19:33:05.060833 | system:policy:conjur
 system:policy:conjur/replication-sets/full                | 2023-06-12 19:33:05.206555 | system:policy:conjur/replication-sets
 system:group:conjur/replication-sets/full/replicated-data | 2023-06-12 19:33:05.206555 | system:policy:conjur/replication-sets
 demo:policy:root                                          | 2023-06-12 19:33:24.699445 |
 demo:group:team_leads                                     | 2023-06-12 19:33:24.699445 | demo:policy:root
(10 rows)

conjur=# select * from permissions limit 10;
 privilege |                        resource_id                        |                            role_id                             |          policy_id           | replication_sets
-----------+-----------------------------------------------------------+----------------------------------------------------------------+------------------------------+------------------
 read      | demo:variable:staging/my-app-1/postgres-database/username | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 execute   | demo:variable:staging/my-app-1/postgres-database/username | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 read      | demo:variable:staging/my-app-1/postgres-database/password | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 execute   | demo:variable:staging/my-app-1/postgres-database/password | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 read      | demo:variable:staging/my-app-1/postgres-database/url      | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 execute   | demo:variable:staging/my-app-1/postgres-database/url      | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 read      | demo:variable:staging/my-app-1/postgres-database/port     | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 execute   | demo:variable:staging/my-app-1/postgres-database/port     | demo:group:staging/my-app-1/postgres-database/secrets-users    | demo:policy:staging/my-app-1 | {}
 update    | demo:variable:staging/my-app-1/postgres-database/username | demo:group:staging/my-app-1/postgres-database/secrets-managers | demo:policy:staging/my-app-1 | {}
 update    | demo:variable:staging/my-app-1/postgres-database/password | demo:group:staging/my-app-1/postgres-database/secrets-managers | demo:policy:staging/my-app-1 | {}
(10 rows)

```




Replication:

- Conjur Roles - defines hosts, users, groups, layers, webservices, policies
- Conjur Resource - sec


## Migrating from Conjur MAML to PML

                              role_id                              |       policy_id
-------------------------------------------------------------------+----------------------------
 !:!:root                                                          |
 demo:user:admin                                                   |
 system:user:admin                                                 |
 system:policy:root                                                |
 system:policy:conjur                                              | system:policy:root
 system:policy:conjur/replication-sets                             | system:policy:conjur
 system:policy:conjur/replication-sets/full                        | system:policy:conjur/replication-sets
 system:group:conjur/replication-sets/full/replicated-data         | system:policy:conjur/replication-sets
 demo:policy:root                                                  |
 demo:group:team_leads                                             | demo:policy:root
 demo:group:security_ops                                           | demo:policy:root
 demo:policy:staging                                               | demo:policy:root
 demo:policy:production                                            | demo:policy:root
 demo:policy:staging/my-app-1                                      | demo:policy:staging
 demo:policy:staging/my-app-2                                      | demo:policy:staging
 demo:policy:staging/my-app-3                                      | demo:policy:staging
 demo:policy:staging/my-app-4                                      | demo:policy:staging
 demo:policy:staging/my-app-5                                      | demo:policy:staging
 demo:policy:staging/my-app-6                                      | demo:policy:staging
 demo:policy:staging/my-app-1/application                          | demo:policy:staging/my-app-1
 demo:layer:staging/my-app-1/application                           | demo:policy:staging/my-app-1
 demo:policy:staging/my-app-1/postgres-database                    | demo:policy:staging/my-app-1
 demo:group:staging/my-app-1/postgres-database/secrets-users       | demo:policy:staging/my-app-1
 demo:group:staging/my-app-1/postgres-database/secrets-managers    | demo:policy:staging/my-app-1
 demo:policy:staging/my-app-2/application                          | demo:policy:staging/my-app-2
 demo:layer:staging/my-app-2/application                           | demo:policy:staging/my-app-2
 demo:policy:staging/my-app-2/postgres-database                    | demo:policy:staging/my-app-2
 demo:group:staging/my-app-2/postgres-database/secrets-users       | demo:policy:staging/my-app-2
 demo:group:staging/my-app-2/postgres-database/secrets-managers    | demo:policy:staging/my-app-2
 demo:policy:staging/my-app-3/application                          | demo:policy:staging/my-app-3
 demo:layer:staging/my-app-3/application                           | demo:policy:staging/my-app-3
 demo:policy:staging/my-app-3/postgres-database                    | demo:policy:staging/my-app-3
 demo:group:staging/my-app-3/postgres-database/secrets-users       | demo:policy:staging/my-app-3
 demo:group:staging/my-app-3/postgres-database/secrets-managers    | demo:policy:staging/my-app-3
 demo:policy:staging/my-app-4/application                          | demo:policy:staging/my-app-4
 demo:layer:staging/my-app-4/application                           | demo:policy:staging/my-app-4
 demo:policy:staging/my-app-4/postgres-database                    | demo:policy:staging/my-app-4
 demo:group:staging/my-app-4/postgres-database/secrets-users       | demo:policy:staging/my-app-4
 demo:group:staging/my-app-4/postgres-database/secrets-managers    | demo:policy:staging/my-app-4
 demo:policy:staging/my-app-5/application                          | demo:policy:staging/my-app-5
 demo:layer:staging/my-app-5/application                           | demo:policy:staging/my-app-5
 demo:policy:staging/my-app-5/postgres-database                    | demo:policy:staging/my-app-5
 demo:group:staging/my-app-5/postgres-database/secrets-users       | demo:policy:staging/my-app-5
 demo:group:staging/my-app-5/postgres-database/secrets-managers    | demo:policy:staging/my-app-5
 demo:policy:staging/my-app-6/application                          | demo:policy:staging/my-app-6
 demo:layer:staging/my-app-6/application                           | demo:policy:staging/my-app-6
 demo:policy:staging/my-app-6/postgres-database                    | demo:policy:staging/my-app-6
 demo:group:staging/my-app-6/postgres-database/secrets-users       | demo:policy:staging/my-app-6
 demo:group:staging/my-app-6/postgres-database/secrets-managers    | demo:policy:staging/my-app-6
 demo:policy:production/my-app-1                                   | demo:policy:production
 demo:policy:production/my-app-2                                   | demo:policy:production
 demo:policy:production/my-app-3                                   | demo:policy:production
 demo:policy:production/my-app-4                                   | demo:policy:production
 demo:policy:production/my-app-5                                   | demo:policy:production
 demo:policy:production/my-app-6                                   | demo:policy:production
 demo:policy:production/my-app-1/application                       | demo:policy:production/my-app-1
 demo:layer:production/my-app-1/application                        | demo:policy:production/my-app-1
 demo:policy:production/my-app-1/postgres-database                 | demo:policy:production/my-app-1
 demo:group:production/my-app-1/postgres-database/secrets-users    | demo:policy:production/my-app-1
 demo:group:production/my-app-1/postgres-database/secrets-managers | demo:policy:production/my-app-1
 demo:policy:production/my-app-2/application                       | demo:policy:production/my-app-2
 demo:layer:production/my-app-2/application                        | demo:policy:production/my-app-2
 demo:policy:production/my-app-2/postgres-database                 | demo:policy:production/my-app-2
 demo:group:production/my-app-2/postgres-database/secrets-users    | demo:policy:production/my-app-2
 demo:group:production/my-app-2/postgres-database/secrets-managers | demo:policy:production/my-app-2
 demo:policy:production/my-app-3/application                       | demo:policy:production/my-app-3
 demo:layer:production/my-app-3/application                        | demo:policy:production/my-app-3
 demo:policy:production/my-app-3/postgres-database                 | demo:policy:production/my-app-3
 demo:group:production/my-app-3/postgres-database/secrets-users    | demo:policy:production/my-app-3
 demo:group:production/my-app-3/postgres-database/secrets-managers | demo:policy:production/my-app-3
 demo:policy:production/my-app-4/application                       | demo:policy:production/my-app-4
 demo:layer:production/my-app-4/application                        | demo:policy:production/my-app-4
 demo:policy:production/my-app-4/postgres-database                 | demo:policy:production/my-app-4
 demo:group:production/my-app-4/postgres-database/secrets-users    | demo:policy:production/my-app-4
 demo:group:production/my-app-4/postgres-database/secrets-managers | demo:policy:production/my-app-4
 demo:policy:production/my-app-5/application                       | demo:policy:production/my-app-5
 demo:layer:production/my-app-5/application                        | demo:policy:production/my-app-5
 demo:policy:production/my-app-5/postgres-database                 | demo:policy:production/my-app-5
 demo:group:production/my-app-5/postgres-database/secrets-users    | demo:policy:production/my-app-5
 demo:group:production/my-app-5/postgres-database/secrets-managers | demo:policy:production/my-app-5
 demo:policy:production/my-app-6/application                       | demo:policy:production/my-app-6
 demo:layer:production/my-app-6/application                        | demo:policy:production/my-app-6
 demo:policy:production/my-app-6/postgres-database                 | demo:policy:production/my-app-6
 demo:group:production/my-app-6/postgres-database/secrets-users    | demo:policy:production/my-app-6
 demo:group:production/my-app-6/postgres-database/secrets-managers | demo:policy:production/my-app-6
