-- Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
--
-- This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
-- found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
-- THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
--
-- The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
-- portions of the software.
--
-- Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
-- SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
-- SOFTWARE.
--
-- Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
-- CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
-- ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  APPLICABLE LAW.


create table if not exists github_installations
(
  id                 uuid,
  team_id            uuid,
  installation_id    bigint,
  installation_type  varchar,
  installation_login varchar,
  sender_id          bigint,
  sender_login       varchar,
  status             varchar,
  created_at         timestamp,
  updated_at         timestamp,
  primary key (id),
) with comment = 'github app installation association with the team';

create index if not exists on github_installations (team_id);
create index if not exists on github_installations (installation_id);
