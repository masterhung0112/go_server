--
-- PostgreSQL database dump
--

-- Dumped from database version 11.7 (Ubuntu 11.7-2.pgdg19.10+1)
-- Dumped by pg_dump version 11.7 (Ubuntu 11.7-2.pgdg19.10+1)

--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: hkuser
--
INSERT INTO public.roles VALUES ('6tt4bj3iztgw7yfe8y97zyoo7o', 'system_post_all', 'authentication.roles.system_post_all.name', 'authentication.roles.system_post_all.description', 1605163387743, 1605163387988, 0, ' create_post use_group_mentions use_channel_mentions', false, true);
INSERT INTO public.roles VALUES ('gkegg9mqi3rgbm9u444mnxkmbc', 'team_post_all_public', 'authentication.roles.team_post_all_public.name', 'authentication.roles.team_post_all_public.description', 1605163387738, 1605163387977, 0, ' use_group_mentions create_post_public use_channel_mentions', false, true);
INSERT INTO public.roles VALUES ('peooyqpsq7g5bfnfo45zb1jiro', 'system_guest', 'authentication.roles.global_guest.name', 'authentication.roles.global_guest.description', 1605163387739, 1605163387978, 0, ' create_direct_channel create_group_channel', true, true);
INSERT INTO public.roles VALUES ('wxat9mo53tg79xdzn55kdq148w', 'channel_admin', 'authentication.roles.channel_admin.name', 'authentication.roles.channel_admin.description', 1605163387741, 1605163387981, 0, ' manage_public_channel_members use_group_mentions add_reaction read_public_channel_groups create_post read_private_channel_groups manage_private_channel_members use_channel_mentions remove_reaction manage_channel_roles', true, true);
INSERT INTO public.roles VALUES ('tkioqq1sgtribqgjbzwop1846c', 'system_read_only_admin', 'authentication.roles.system_read_only_admin.name', 'authentication.roles.system_read_only_admin.description', 1605163387745, 1605163387990, 0, ' list_private_teams sysconsole_read_user_management_users sysconsole_read_experimental read_other_users_teams sysconsole_read_environment sysconsole_read_user_management_channels read_public_channel sysconsole_read_user_management_permissions read_jobs read_public_channel_groups sysconsole_read_integrations list_public_teams sysconsole_read_user_management_groups sysconsole_read_about sysconsole_read_reporting sysconsole_read_site read_private_channel_groups read_channel sysconsole_read_authentication view_team sysconsole_read_plugins sysconsole_read_user_management_teams', false, true);
INSERT INTO public.roles VALUES ('mrejpofuoffiiynqcsi98es9ya', 'channel_guest', 'authentication.roles.channel_guest.name', 'authentication.roles.channel_guest.description', 1605163387746, 1605163387991, 0, ' create_post use_channel_mentions use_slash_commands read_channel add_reaction remove_reaction upload_file edit_post', true, true);
INSERT INTO public.roles VALUES ('96whs8mg73dszp7cz4u7sdbd7c', 'team_guest', 'authentication.roles.team_guest.name', 'authentication.roles.team_guest.description', 1605163387741, 1605163387982, 0, ' view_team', true, true);
INSERT INTO public.roles VALUES ('13kpq8iaqffmdf9qkrfqmpby9h', 'team_admin', 'authentication.roles.team_admin.name', 'authentication.roles.team_admin.description', 1605163387742, 1605163387983, 0, ' import_team convert_private_channel_to_public read_private_channel_groups manage_team convert_public_channel_to_private manage_channel_roles use_group_mentions manage_private_channel_members remove_reaction manage_others_outgoing_webhooks delete_others_posts manage_incoming_webhooks remove_user_from_team add_reaction manage_outgoing_webhooks manage_slash_commands manage_public_channel_members delete_post manage_others_slash_commands use_channel_mentions manage_team_roles create_post read_public_channel_groups manage_others_incoming_webhooks', true, true);
INSERT INTO public.roles VALUES ('7ta1wfbacjy3zxid54n3cqjzqw', 'system_post_all_public', 'authentication.roles.system_post_all_public.name', 'authentication.roles.system_post_all_public.description', 1605163387742, 1605163387985, 0, ' use_group_mentions create_post_public use_channel_mentions', false, true);
INSERT INTO public.roles VALUES ('rfc1w7z71pnzurkhpb1jgrbmdh', 'team_user', 'authentication.roles.team_user.name', 'authentication.roles.team_user.description', 1605163387747, 1605163387992, 0, ' list_team_channels join_public_channels read_public_channel view_team create_public_channel create_private_channel invite_user add_user_to_team', true, true);
INSERT INTO public.roles VALUES ('nh5i9ik1u78hdcny9usdoixkuo', 'channel_user', 'authentication.roles.channel_user.name', 'authentication.roles.channel_user.description', 1605163387735, 1605163387974, 0, ' manage_private_channel_members manage_private_channel_properties get_public_link manage_public_channel_properties use_slash_commands remove_reaction create_post add_reaction delete_post use_channel_mentions edit_post delete_private_channel upload_file delete_public_channel read_public_channel_groups manage_public_channel_members read_private_channel_groups use_group_mentions read_channel', true, true);
INSERT INTO public.roles VALUES ('tj3atgnwjfrt7emz8pgqmh5z4c', 'team_post_all', 'authentication.roles.team_post_all.name', 'authentication.roles.team_post_all.description', 1605163387737, 1605163387975, 0, ' create_post use_channel_mentions use_group_mentions', false, true);
INSERT INTO public.roles VALUES ('xf95ytghtjfsfd543dum68uzua', 'system_user_access_token', 'authentication.roles.system_user_access_token.name', 'authentication.roles.system_user_access_token.description', 1605163387743, 1605163387986, 0, ' revoke_user_access_token create_user_access_token read_user_access_token', false, true);
INSERT INTO public.roles VALUES ('hm1bxei8b3d68e4j95tqnndppw', 'system_manager', 'authentication.roles.system_manager.name', 'authentication.roles.system_manager.description', 1605163387740, 1605163387980, 0, ' sysconsole_write_user_management_groups sysconsole_read_user_management_teams read_jobs manage_private_channel_members convert_private_channel_to_public read_public_channel sysconsole_write_user_management_permissions sysconsole_write_environment manage_channel_roles sysconsole_write_user_management_channels manage_public_channel_properties read_private_channel_groups add_user_to_team sysconsole_read_about manage_team_roles view_team edit_brand sysconsole_write_user_management_teams join_private_teams manage_team sysconsole_read_site sysconsole_read_user_management_permissions sysconsole_read_reporting list_private_teams manage_jobs sysconsole_write_site read_public_channel_groups list_public_teams delete_public_channel sysconsole_read_environment sysconsole_read_authentication sysconsole_read_user_management_groups sysconsole_read_plugins delete_private_channel sysconsole_write_integrations sysconsole_read_user_management_channels convert_public_channel_to_private manage_public_channel_members read_channel manage_private_channel_properties remove_user_from_team sysconsole_read_integrations join_public_teams', false, true);
INSERT INTO public.roles VALUES ('f9drbz6cyjdmb8jof6smiqya7h', 'system_user_manager', 'authentication.roles.system_user_manager.name', 'authentication.roles.system_user_manager.description', 1605163387744, 1605163387989, 0, ' add_user_to_team manage_private_channel_properties sysconsole_read_user_management_groups read_public_channel convert_public_channel_to_private read_public_channel_groups read_jobs read_private_channel_groups list_private_teams manage_public_channel_properties join_private_teams sysconsole_read_user_management_permissions manage_public_channel_members read_channel sysconsole_read_user_management_teams delete_private_channel sysconsole_read_user_management_channels list_public_teams manage_private_channel_members join_public_teams manage_team remove_user_from_team sysconsole_write_user_management_channels sysconsole_write_user_management_groups sysconsole_write_user_management_teams convert_private_channel_to_public view_team delete_public_channel manage_channel_roles manage_team_roles sysconsole_read_authentication', false, true);
INSERT INTO public.roles VALUES ('d54xjt4sat8h7dqwu6i35jocuy', 'system_user', 'authentication.roles.global_user.name', 'authentication.roles.global_user.description', 1605163387739, 1605163387993, 0, ' list_public_teams join_public_teams create_direct_channel create_group_channel view_members create_team create_emojis delete_emojis', true, true);
INSERT INTO public.roles VALUES ('ha8u9qxwx3dm8mnbq8sfi7ugdc', 'system_admin', 'authentication.roles.global_admin.name', 'authentication.roles.global_admin.description', 1605163387745, 1605163387995, 0, ' sysconsole_read_site manage_jobs sysconsole_write_user_management_channels read_public_channel_groups sysconsole_read_user_management_permissions sysconsole_read_environment manage_others_outgoing_webhooks manage_outgoing_webhooks sysconsole_write_plugins sysconsole_write_user_management_teams list_private_teams import_team read_other_users_teams create_post_ephemeral read_user_access_token read_public_channel create_team get_public_link create_emojis delete_post manage_public_channel_properties delete_others_posts read_jobs manage_others_incoming_webhooks create_public_channel use_slash_commands sysconsole_read_experimental invite_user sysconsole_write_site manage_others_slash_commands list_public_teams assign_bot read_channel convert_public_channel_to_private sysconsole_write_user_management_permissions read_private_channel_groups manage_channel_roles edit_post create_user_access_token sysconsole_read_user_management_users edit_brand sysconsole_read_about list_users_without_team read_bots manage_private_channel_members create_group_channel delete_others_emojis manage_team sysconsole_write_user_management_users upload_file create_post manage_slash_commands sysconsole_write_about list_team_channels create_private_channel sysconsole_write_environment read_others_bots sysconsole_write_authentication manage_bots delete_private_channel join_public_teams manage_shared_channels create_post_public use_channel_mentions edit_other_users manage_incoming_webhooks join_private_teams sysconsole_write_compliance manage_system manage_others_bots sysconsole_read_user_management_groups view_team sysconsole_read_compliance add_user_to_team sysconsole_read_integrations sysconsole_write_user_management_groups sysconsole_write_experimental manage_team_roles join_public_channels manage_private_channel_properties manage_roles promote_guest invite_guest convert_private_channel_to_public sysconsole_write_reporting assign_system_admin_role revoke_user_access_token remove_user_from_team sysconsole_read_user_management_channels sysconsole_read_plugins remove_reaction add_reaction delete_public_channel view_members edit_others_posts sysconsole_write_integrations sysconsole_read_user_management_teams delete_emojis sysconsole_read_authentication create_direct_channel create_bot sysconsole_read_reporting use_group_mentions demote_to_guest remove_others_reactions manage_oauth manage_system_wide_oauth manage_public_channel_members', true, true);

--
-- Data for Name: systems; Type: TABLE DATA; Schema: public; Owner: hkuser
--
INSERT INTO public.systems VALUES ('emoji_permissions_split', 'true');
INSERT INTO public.systems VALUES ('webhook_permissions_split', 'true');
INSERT INTO public.systems VALUES ('list_join_public_private_teams', 'true');
INSERT INTO public.systems VALUES ('remove_permanent_delete_user', 'true');
INSERT INTO public.systems VALUES ('add_bot_permissions', 'true');
INSERT INTO public.systems VALUES ('apply_channel_manage_delete_to_channel_user', 'true');
INSERT INTO public.systems VALUES ('remove_channel_manage_delete_from_team_user', 'true');
INSERT INTO public.systems VALUES ('view_members_new_permission', 'true');
INSERT INTO public.systems VALUES ('add_manage_guests_permissions', 'true');
INSERT INTO public.systems VALUES ('channel_moderations_permissions', 'true');
INSERT INTO public.systems VALUES ('add_use_group_mentions_permission', 'true');
INSERT INTO public.systems VALUES ('add_system_console_permissions', 'true');
INSERT INTO public.systems VALUES ('add_convert_channel_permissions', 'true');
INSERT INTO public.systems VALUES ('manage_shared_channel_permissions', 'true');

--
-- PostgreSQL database dump complete
--
