-- Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
--
-- This software is made available by Breu, Inc., under the terms of the Breu  
-- Community License Agreement, Version 1.0 located at  
-- http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
-- ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
-- OF SUCH LICENSE AGREEMENT. 

 

create table if not exists github_repos (
  id                uuid,
  team_id           uuid,
  github_id         bigint,
  name              varchar,
  full_name         varchar,
  created_at        timestamp,
  updated_at        timestamp,
  primary key       (id),
) with comment = 'repositories provided by github installation';

create index if not exists on github_repos (team_id);
create index if not exists on github_repos (github_id);

