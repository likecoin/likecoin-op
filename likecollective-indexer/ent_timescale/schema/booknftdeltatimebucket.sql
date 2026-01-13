-- 7d

CREATE MATERIALIZED VIEW book_nft_delta_time_bucket_7d
WITH (timescaledb.continuous) AS
SELECT 
    nft_class_address || '_' || time_bucket(INTERVAL '7 day', datetime) as id,
    nft_class_address as evm_address,
    time_bucket(INTERVAL '7 day', datetime) AS bucket,
    sum(staked_amount_added) as staked_amount,
    max(datetime) as last_staked_at,
    count(distinct account_evm_address) as number_of_stakers
FROM staking_events_hyper_table
WHERE event_type = 'staked'
GROUP BY nft_class_address, bucket;

-- Use real-time aggregates
ALTER MATERIALIZED VIEW book_nft_delta_time_bucket_7d set (timescaledb.materialized_only = false);

SELECT add_continuous_aggregate_policy('book_nft_delta_time_bucket_7d',
  start_offset => INTERVAL '21 day',
  end_offset => NULL,
  schedule_interval => INTERVAL '1 day');

-- end 7d

-- 30d

CREATE MATERIALIZED VIEW book_nft_delta_time_bucket_30d
WITH (timescaledb.continuous) AS
SELECT
    nft_class_address || '_' || time_bucket(INTERVAL '30 day', datetime) as id,
    nft_class_address as evm_address,
    time_bucket(INTERVAL '30 day', datetime) AS bucket,
    sum(staked_amount_added) as staked_amount,
    max(datetime) as last_staked_at,
    count(distinct account_evm_address) as number_of_stakers
FROM staking_events_hyper_table
WHERE event_type = 'staked'
GROUP BY nft_class_address, bucket;

-- Use real-time aggregates
ALTER MATERIALIZED VIEW book_nft_delta_time_bucket_30d set (timescaledb.materialized_only = false);

SELECT add_continuous_aggregate_policy('book_nft_delta_time_bucket_30d',
  start_offset => INTERVAL '90 day',
  end_offset => NULL,
  schedule_interval => INTERVAL '1 day');

-- end 30d

-- 1y

CREATE MATERIALIZED VIEW book_nft_delta_time_bucket_1y
WITH (timescaledb.continuous) AS
SELECT
    nft_class_address || '_' || time_bucket(INTERVAL '1 year', datetime) as id,
    nft_class_address as evm_address,
    time_bucket(INTERVAL '1 year', datetime) AS bucket,
    sum(staked_amount_added) as staked_amount,
    max(datetime) as last_staked_at,
    count(distinct account_evm_address) as number_of_stakers
FROM staking_events_hyper_table
WHERE event_type = 'staked'
GROUP BY nft_class_address, bucket;

-- Use real-time aggregates
ALTER MATERIALIZED VIEW book_nft_delta_time_bucket_1y set (timescaledb.materialized_only = false);

SELECT add_continuous_aggregate_policy('book_nft_delta_time_bucket_1y',
  start_offset => INTERVAL '3 year',
  end_offset => NULL,
  schedule_interval => INTERVAL '1 day');

-- end 1y
