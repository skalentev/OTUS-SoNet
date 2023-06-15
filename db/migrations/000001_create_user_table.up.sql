CREATE TABLE IF NOT EXISTS public."user" (
                               id character varying(64)  NOT NULL,
                               created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                               deleted_at timestamp without time zone,
                               first_name character varying(90) NOT NULL,
                               second_name character varying(90) NOT NULL,
                               birthdate character varying(20) NOT NULL,
                               biography text,
                               city character varying(64),
                               password character varying(64)
);
ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);
CREATE INDEX IF NOT EXISTS idx_first_name ON public."user" USING btree (first_name);
CREATE INDEX IF NOT EXISTS idx_second_name ON public."user" USING btree (second_name);
CREATE TABLE IF NOT EXISTS public.session (
                                token character varying(64) NOT NULL,
                                user_id character varying(64) NOT NULL,
                                created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                deleted_at timestamp without time zone,
                                token_till timestamp without time zone
);
ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(id);
CREATE TABLE IF NOT EXISTS public.friend
(
    id character varying(64)  NOT NULL,
    user_id character varying(64)  NOT NULL,
    friend_id character varying(64)  NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone
);
ALTER TABLE IF EXISTS public.friend
    OWNER to "user";
CREATE INDEX IF NOT EXISTS idx_friend_friend_id
    ON public.friend USING btree
        (friend_id  ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_friend_user_id
    ON public.friend USING btree
        (user_id  ASC NULLS LAST);
CREATE TABLE IF NOT EXISTS public.post
(
    id character varying(64)  NOT NULL,
    user_id character varying(64)  NOT NULL,
    text text ,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    CONSTRAINT pk_post_id PRIMARY KEY (id)
);
ALTER TABLE IF EXISTS public.post
    OWNER to "user";
CREATE INDEX IF NOT EXISTS idx_post_created_at
    ON public.post USING btree
        (created_at ASC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_post_user_id
    ON public.post USING btree
        (user_id ASC NULLS LAST);
