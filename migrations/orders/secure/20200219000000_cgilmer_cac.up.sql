-- This migration allows a CAC cert to have read/write access to all orders.
-- The Orders API uses client certificate authentication. Only certificates
-- signed by a trusted CA (such as DISA) are allowed which includes CACs.
-- Using a person's CAC as the certificate is a convenient way to permit a
-- single trusted individual to interact with the Orders API and the Prime API. Eventually
-- this CAC certificate should be removed.
INSERT INTO public.client_certs (
	id,
	sha256_digest,
	subject,
	allow_orders_api,
	allow_air_force_orders_read,
	allow_air_force_orders_write,
	allow_army_orders_read,
	allow_army_orders_write,
	allow_coast_guard_orders_read,
	allow_coast_guard_orders_write,
	allow_marine_corps_orders_read,
	allow_marine_corps_orders_write,
	allow_navy_orders_read,
	allow_navy_orders_write,
	created_at,
	updated_at)
VALUES (
	'01e5d7bb-b4b4-41c7-8101-7c0a23b8aa17',
	'15697aff85d93a2f04259618b66f09f9976887b3f2ee206395c42e4897aefe61',
	'CN=chrisgilmerproj,OU=DoD+OU=PKI+OU=CONTRACTOR,O=U.S. Government,C=US',
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	true,
	now(),
	now());
