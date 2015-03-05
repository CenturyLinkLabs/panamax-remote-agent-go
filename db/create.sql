CREATE TABLE "deployments"
(
"id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"name" varchar(255) NOT NULL,
"template" text NOT NULL,
"service_ids" text
);
