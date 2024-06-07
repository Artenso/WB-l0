-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders(
    order_uid TEXT PRIMARY KEY,
    order_json JSON
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
