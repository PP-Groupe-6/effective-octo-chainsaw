DROP TABLE IF EXISTS "account" CASCADE;
CREATE TABLE "account" (
	"client_id" varchar(255) NOT NULL,
	"name" varchar(40) NOT NULL,
	"surname" varchar(40) NOT NULL,
	"phone_number" varchar(40) NOT NULL,
	"mail_adress" varchar(40) NOT NULL,
	"account_amount" numeric NOT NULL,
	CONSTRAINT "Client_pk" PRIMARY KEY ("client_id")
) WITH (
  OIDS=FALSE
);


DROP TABLE IF EXISTS "transfer";
CREATE TABLE "transfer" (
	"transfer_id" varchar(255) NOT NULL,
	"transfer_type" varchar(255) NOT NULL,
	"transfer_amount" numeric NOT NULL,
	"account_transfer_payer_id" varchar(255) NOT NULL,
	"account_transfer_receiver_id" varchar(255) NOT NULL,
	"receiver_question" varchar(255) NOT NULL,
	"receiver_answer" varchar(255) NOT NULL,
	"scheduled_transfer_date" DATE,
	CONSTRAINT "Transfer_pk" PRIMARY KEY ("transfer_id")
) WITH (
  OIDS=FALSE
);


DROP TABLE IF EXISTS "invoice";
CREATE TABLE "invoice" (
	"invoice_id" varchar(255) NOT NULL,
	"invoice_amount" numeric NOT NULL,
	"invoice_state" int NOT NULL,
	"invoice_expiration_date" DATE NOT NULL,
	"account_invoice_payer_id" varchar(255) NOT NULL,
	"account_invoice_receiver_id" varchar(255) NOT NULL,
	CONSTRAINT "Invoice_pk" PRIMARY KEY ("invoice_id")
) WITH (
  OIDS=FALSE
);


DROP TABLE IF EXISTS "invoice_state";
CREATE TABLE "invoice_state" (
	"state_id" int NOT NULL,
	"state_name" varchar(40) NOT NULL,
	CONSTRAINT "InvoiceState_pk" PRIMARY KEY ("state_id")
) WITH (
  OIDS=FALSE
);

ALTER TABLE "transfer" ADD CONSTRAINT "Transfer_fk0" FOREIGN KEY ("account_transfer_payer_id") REFERENCES "account"("client_id");
ALTER TABLE "transfer" ADD CONSTRAINT "Transfer_fk1" FOREIGN KEY ("account_transfer_receiver_id") REFERENCES "account"("client_id");

ALTER TABLE "invoice" ADD CONSTRAINT "Invoice_fk0" FOREIGN KEY ("invoice_state") REFERENCES "invoice_state"("state_id");
ALTER TABLE "invoice" ADD CONSTRAINT "Invoice_fk1" FOREIGN KEY ("account_invoice_payer_id") REFERENCES "account"("client_id");
ALTER TABLE "invoice" ADD CONSTRAINT "Invoice_fk2" FOREIGN KEY ("account_invoice_receiver_id") REFERENCES "account"("client_id");