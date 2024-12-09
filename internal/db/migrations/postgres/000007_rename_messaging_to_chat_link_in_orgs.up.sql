-- Step 1: Alter the table to update the default value of the 'hooks' column
ALTER TABLE orgs
ALTER COLUMN hooks SET DEFAULT '{"repo":0, "chat_links":0}';