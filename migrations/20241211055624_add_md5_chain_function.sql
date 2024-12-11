-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION md5_concat(text, text) returns text as
 'select md5($1 || $2);' language sql;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE AGGREGATE md5_chain(text) (
    sfunc = md5_concat,
    stype = text,
    initcond = ''
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP AGGREGATE md5_chain(text);
-- +goose StatementEnd
-- +goose StatementBegin
DROP FUNCTION md5_concat(text, text);
-- +goose StatementEnd
