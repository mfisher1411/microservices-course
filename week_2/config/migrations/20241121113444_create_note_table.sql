-- +goose Up
    create table note (
        id serial primary key,
        title text not null,
        body text not null,
        created_at timestamp not null default now(),
        updated_at timestamp
    );
-- +goose StatementBegin
-- SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
    drop table note;
-- +goose StatementBegin
-- SELECT 'down SQL query';
-- +goose StatementEnd
