CREATE TABLE IF NOT EXISTS `page_widgets` (
	`id_widget`	INTEGER NOT NULL,
	`id_page`	INTEGER NOT NULL,
	PRIMARY KEY(id_widget,id_page)
  FOREIGN KEY (id_widget) REFERENCES widgets(id)
	FOREIGN KEY (id_page) REFERENCES pages(id)
);
