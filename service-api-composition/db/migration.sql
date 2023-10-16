create table if not exists scale_result
(
    task_id         text primary key,
    origin_image_id text not null,
    scale_factor    int  not null,

    image_id        text,
    error           text,

    constraint image_id_error_both_not_null check (image_id is not null or error is not null)
);
