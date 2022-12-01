-- Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
--
-- This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
-- found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
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
-- ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.


create table if not exists workloads (
  	id            uuid,
		app_id        uuid,
		repo_id       uuid,
		repo_path     varchar,
		name          varchar,
		kind          varchar,
		container     varchar,
		builder       varchar,
		created_at    timestamp,
		updated_at    timestamp,
    primary key   (id),
) with comment = 'workload defines a single unit of compute';

create table if not exists resources (
  	id            uuid,
		app_id        uuid,
		name          varchar,
		provider      varchar,
		driver        varchar,
		is_immutable  boolean,
		created_at    timestamp,
		updated_at    timestamp,
    primary key   (id),
) with comment = 'resource is cloud provider specific';


create table if not exists blueprints (
    id              uuid,
		app_id          uuid,
		repo_id         uuid,
		name            varchar,
		regions         varchar,
		repo_branch     varchar,
		rollout_budget  varchar,
		created_at      timestamp,
		updated_at      timestamp,
    primary key     (id),
) with comment = 'blueprint defines the deployment and how to rollout';

create table if not exists rollouts (
  	id                  uuid,
		app_id              uuid,
		blueprint_id        uuid,
		trigger             varchar,
		state               varchar,
		created_at          timestamp,
		updated_at          timestamp,
    primary key         (id),
) with comment = 'rollouts is the result of a blueprint being applied';
