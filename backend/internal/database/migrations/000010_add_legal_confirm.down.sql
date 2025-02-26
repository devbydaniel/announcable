-- Drop indexes for privacy_policy_confirm table
DROP INDEX IF EXISTS idx_privacy_policy_confirm_user_version;
DROP INDEX IF EXISTS idx_privacy_policy_confirm_version;
DROP INDEX IF EXISTS idx_privacy_policy_confirm_user_id;

-- Drop privacy_policy_confirm table
DROP TABLE IF EXISTS privacy_policy_confirms;

-- Drop indexes for tos_confirm table
DROP INDEX IF EXISTS idx_tos_confirm_user_version;
DROP INDEX IF EXISTS idx_tos_confirm_version;
DROP INDEX IF EXISTS idx_tos_confirm_user_id;

-- Drop tos_confirm table
DROP TABLE IF EXISTS tos_confirms;
