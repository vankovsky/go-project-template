package migrations

type Migration struct {
	Name       string
	Statements []string
}

var (
// 	AddAgent = Migration{
// 		Name: "add host table",
// 		Statements: []string{
// 			`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,
// 			`create table if not exists public.agent
// (
//     "id" uuid NOT NULL DEFAULT gen_random_uuid(),
//     host      varchar not null unique,
//     ip        inet    not null,
//     name  varchar not null,
//     meta      jsonb,
//     is_active boolean default false,
//     PRIMARY KEY ("id"),
//     unique (id, host)
// )`,

// 			`comment on table public.agent is 'A list of all registered agents'`,
// 			`comment on column public.agent.host is 'A host where the agents is running'`,
// 			`comment on column public.agent.ip is 'A host''s IP  where the agents is running'`,
// 			`comment on column public.agent.name is 'An agent name (e.g. Moscow, Selectel)'`,
// 			`comment on column public.agent.meta is 'An extended description data for client apps (e.g. location, title in different languages)'`,
// 			`comment on column public.agent.is_active is 'Whether an agent is active and can perform their tasks'`,
// 			`insert into public.agent (id, host, ip, name, meta, is_active)
// 			values (DEFAULT, 'localhost', '127.0.0.1/32'::inet, 'local', null, true)
// 			on conflict do nothing`,
// 		},
// 	}

// 	AddHost = Migration{
// 		Name: "add host table",
// 		Statements: []string{
// 			`create table if not exists public.host
// (
//     "id" uuid NOT NULL DEFAULT gen_random_uuid(),
//     name varchar not null,
//     ips jsonb,
//     PRIMARY KEY ("id"),
//     unique (id, name)
// )`,
// 			`comment on table public.host is 'A list of host names, which are derived from the URLs'`,
// 			`comment on column public.host.name is 'A host name derived from a URL'`,
// 			`comment on column public.host.ips is 'a JSON of IPs of the host (e.g {"a":[1.2.3.4, 5.6.7.8], "mx":[{"ip":"1.2.3.4", "prio":5}}]})'`,
// 		},
// 	}

// 	AddURL = Migration{
// 		Name: "add url table",
// 		Statements: []string{
// 			`create table if not exists public.url
// (
//     "id" uuid NOT NULL DEFAULT gen_random_uuid(),
//     url        varchar not null,
//     is_active  boolean                  default true,
//     created_at timestamp with time zone default CURRENT_TIMESTAMP,
//     updated_at timestamp with time zone default CURRENT_TIMESTAMP,
//     host_id    uuid not null constraint url_host_id_fk references host on update cascade on delete cascade,
//     settings   jsonb,
//     PRIMARY KEY ("id"),
//     unique (url)
// )`,
// 			`comment on table public.url is 'A list of URLs to check'`,
// 			`comment on column public.url.url is 'A unique URL to check'`,
// 			`comment on column public.url.is_active is 'Whether to check the URL'`,
// 			`comment on column public.url.host_id is 'A reference to a derived from URL host'`,
// 			`comment on column public.url.settings is 'A meta info for the URL (e.g. auth. creds., headers)'`,
// 		},
// 	}

// 	AddURLAgent = Migration{
// 		Name: "add url_agent table",
// 		Statements: []string{
// 			`create table if not exists public.url_agent
// (
//     url_id   uuid not null constraint url_agent_url_id_fkey references url (id) on update cascade on delete cascade,
//     agent_id uuid not null
//         constraint url_agent_agent_id_fkey
//             references agent
//             on update cascade,
//     constraint url_agent_pkey
//         primary key (url_id, agent_id)
// )`,
// 		},
// 	}

// 	AddScheduler = Migration{
// 		Name: "add changelog DNS",
// 		Statements: []string{`
// CREATE TABLE if not exists "public"."changelog_dns" (
//     "id" uuid NOT NULL DEFAULT gen_random_uuid(),
//     "host_id" uuid NOT NULL,
//     "agent_id" uuid NOT NULL,
//     "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
//     "ips" jsonb,
//     "name_servers" jsonb,
//     "error" jsonb,
//     "do_err_alert" bool,
//     CONSTRAINT "changelog_dns_host_id_fk" FOREIGN KEY ("host_id") REFERENCES "public"."host"("id") ON DELETE CASCADE ON UPDATE CASCADE,
//     CONSTRAINT "changelog_dns_agent_id_fk" FOREIGN KEY ("agent_id") REFERENCES "public"."agent"("id") ON DELETE CASCADE ON UPDATE CASCADE,
//     PRIMARY KEY ("id")
// )`,
// 		},
// 	}

// 	AddCheckDNS = Migration{
// 		Name: "add checks log",
// 		Statements: []string{
// 			`CREATE TABLE if not exists "public"."checks_log" (
//     "check_date" date NOT NULL,
//     "checked_at" timestamptz NOT NULL,
//     "agent_id" uuid NOT NULL,
//     "result" jsonb NOT NULL,
//     "host_id" uuid NOT NULL,
//     "check_type" int2 NOT NULL,
//     CONSTRAINT "checks_dns_agent_id_fk" FOREIGN KEY ("agent_id") REFERENCES "public"."agent"("id") ON DELETE CASCADE ON UPDATE CASCADE,
//     CONSTRAINT "checks_dns_host_id_fk" FOREIGN KEY ("host_id") REFERENCES "public"."host"("id") ON DELETE CASCADE ON UPDATE CASCADE
// ) partition BY RANGE (check_date)`,
// 			`COMMENT ON COLUMN "public"."checks_log"."checked_at" IS 'when the check has been completed'`,
// 			`COMMENT ON COLUMN "public"."checks_log"."agent_id" IS 'which agent performed the ckeck'`,
// 			`COMMENT ON COLUMN "public"."checks_log"."result" IS 'a result of DNS check (e.g A, MX records or error)'`,
// 			`COMMENT ON COLUMN "public"."checks_log"."host_id" IS 'against which host the check was performed'`,
// 			`COMMENT ON TABLE "public"."checks_log" IS 'a log with performed checks'`,
// 		},
// 	}

// 	AddSchduler2 = Migration{
// 		Name: "add schedule table",
// 		Statements: []string{
// 			`create table if not exists public.schedule
// (
//     entity_type smallint,
//     payload     jsonb,
//     check_at    timestamp with time zone not null,
//     status      smallint,
//     agent_id    integer constraint schedule_agent_id_fkey references agent on update cascade on delete cascade,
//     url_id      integer constraint schedule_url_id_fk references url (id) on update cascade on delete cascade,
//     host_id     integer constraint schedule_host_id_fk references host on update cascade

// )`}}
)
