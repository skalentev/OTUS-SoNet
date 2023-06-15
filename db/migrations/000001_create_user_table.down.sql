DROP TABLE IF EXISTS public.user;
DROP TABLE IF EXISTS public.session;
DROP TABLE IF EXISTS public.friend;
DROP INDEX IF EXISTS public.idx_friend_friend_id;
DROP INDEX IF EXISTS public.idx_friend_user_id;
DROP TABLE IF EXISTS public.post;
DROP INDEX IF EXISTS public.idx_post_created_at;
DROP INDEX IF EXISTS public.idx_post_user_id;
