CREATE TABLE electronic_orders_revisions (
	id uuid PRIMARY KEY,
	electronic_order_id uuid NOT NULL,
	FOREIGN KEY (electronic_order_id) REFERENCES electronic_orders (id) ON DELETE CASCADE,
	seq_num INT DEFAULT 0 NOT NULL,
	given_name text NOT NULL,
	middle_name text NULL,
	family_name text NOT NULL,
	name_suffix text NULL,
	affiliation text NOT NULL,
	paygrade text NOT NULL,
	title text NULL,
	status text NOT NULL,
	date_issued timestamp NOT NULL,
	no_cost_move boolean DEFAULT false NOT NULL,
	tdy_en_route boolean DEFAULT false NOT NULL,
	tour_type text DEFAULT 'accompanied' NOT NULL,
	orders_type text NOT NULL,
	has_dependents boolean NOT NULL,
	losing_uic text NULL,
	losing_unit_name text NULL,
	losing_unit_city text NULL,
	losing_unit_locality text NULL,
	losing_unit_country text NULL,
	losing_unit_postal_code text NULL,
	gaining_uic text NULL,
	gaining_unit_name text NULL,
	gaining_unit_city text NULL,
	gaining_unit_locality text NULL,
	gaining_unit_country text NULL,
	gaining_unit_postal_code text NULL,
	report_no_earlier_than timestamp NULL,
	report_no_later_than timestamp NULL,
	hhg_tac text NULL,
	hhg_sdn text NULL,
	hhg_loa text NULL,
	nts_tac text NULL,
	nts_sdn text NULL,
	nts_loa text NULL,
	pov_shipment_tac text NULL,
	pov_shipment_sdn text NULL,
	pov_shipment_loa text NULL,
	pov_storage_tac text NULL,
	pov_storage_sdn text NULL,
	pov_storage_loa text NULL,
	ub_tac text NULL,
	ub_sdn text NULL,
	ub_loa text NULL,
	comments text NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

CREATE INDEX ON electronic_orders_revisions (electronic_order_id, seq_num);
