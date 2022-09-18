-- Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.  

create keyspace if not exists ctrlplane 
  with replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
