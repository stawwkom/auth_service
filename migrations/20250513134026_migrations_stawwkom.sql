-- +goose Up
-- +goose StatementBegin
create table auth (
    id serial primary key,
    login text not null,
    password text not null,
    context text not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists auth;
-- +goose StatementEnd
