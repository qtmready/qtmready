insert into orgs (id, name, domain, slug)
values ('00000000-0000-0000-0000-000000000001', 'no org', 'example.com', 'no-org')
on conflict (id) do nothing
returning id, created_at, updated_at, name, domain, slug;
