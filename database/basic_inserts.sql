
INSERT INTO compositions (composition_id, energetics_id, taurine, caffeine)
VALUES 
    (1, null, 400, 32), 
    (2, null, 380, 30),
	(3,null,240,30),
	(4,null,40,30),
	(5,null,30,27),
	(6,null,400,32),
	(7,null,400,30),
	(8,null,240,32),
	(9,null,400,30),
	(10,null,50,30),
	(11,null,0,32),
	(12,null,200,50);
	

INSERT INTO energetics (name, taste, description, manufacturer_name, manufacture_country, composition_id, picture_url)
VALUES 
    ('Monster&reg', 'Pipeline Punch', 'A flavour like sakura, glowy barby-like pink metal casing and monstrous notorious but beloved Monster&reg logo on top!', 'Monster Beverage', 'USA', 1, '/static/img/cover_monster_punch_pipeline.png'),
    ('KAIF&reg', 'Mojito Virgin', 'Astonishing and refreshing! Top german bewerage with real mint on the tip of your tongue!', 'Ancor Group Germany GmbH', 'Germany', 2, '/static/img/cover_kaif_mojito.png'),
	('Lit Energy','Blueberry','Peachy soft and wet blueberries from taiga itself were magically transformed into this delicious drink!','Lit Energy','Russia',3,'/static/img/cover_lit_blueberry.png'),
	('Tornado&reg','Cactus','Have you ever tasted a CACTUS? A spiky flavor made of exotic fruits (yes, cactus is a fruit!) that will tingle you up.','Global Functional Drinks Rus Ltd.','Russia',4,'/static/img/cover_tornado_cactus.png'),
	('Gorilla','Mint','Mint tic tac in a liquid form. Don''t freeze your lips!','Gorilla Drinks Limited','Cyprus',5,'/static/img/cover_gorilla_mint.png'),
	('RedBull&reg','Red Edition','Red stands for Really ExiteD4theWATERMELONTASTE! Classical RedBull&reg quality, now flavored with your favourite summer dish!','Red Bull GmbH','Austria',6,'/static/img/cover_redbull_red.png'),
	('NonStop','Original','A classic taste straight from our beloved post-apocalyptic game! A blowout is coming, stalkers...','New Products Group','Ukraine',7,'/static/img/cover_nonstop_original.png'),
	('Carabao&reg','Green Apple','Classic Carabao&reg, now with gentle apple touch. Totally green and totally energizing.','Carabao Tawandang Co.','Thailand',8,'/static/img/cover_carabao_greenapple.png'),
	('Monster&reg','Pacific Punch','Set sail with Monster&reg and this otherworldly exotic citrus flavour with a pinch of cinnamon aroma. Reminds of navy lore...','Monster Beverage','USA',9,'/static/img/cover_monster_punch_pacific.png'),
	('Jaguar&reg','Cult','Forest berries are freshly transformed into iconic 2007 hype drink. Do not overhype, but overcharge!','United Bottling Group','Russia',10,'/static/img/cover_jaguar_cult.png'),
	('CocaCola&reg','Energy','Yes, this thing exists. Now you can throw away all other ''cola'' energy replicas and experience the real brand and original coke flavour.','The Coca-Cola Company','USA',11,'/static/img/cover_cocacola_energy.png'),
	('Firegin','Original','A heartrate killer from Korea. It''s not a Cocaine Energy, but still will punch your grumpy face up!','Firegin Co Ltd','South Korea',12,'/static/img/cover_firegin.png', 
	'Zestful', 'Citrus', 'Energizing burst of tanginess', 'EnerGize Corp.', 'USA',
'Vivid', 'Berry', 'Bold and invigorating flavor explosion', 'VigorMax Inc.', 'Canada',
'Revitalize', 'Tropical', 'Refreshingly exotic with a revitalizing kick', 'RevivePlus Ltd.', 'Australia',
'PowerPunch', 'Pineapple', 'Punchy fusion of tropical sweetness', 'PowerUp Enterprises', 'UK',

'VigorVerve', 'Cranberry', 'Vibrant cranberry kick for a burst of verve', 'VigorVerve Wellness', 'USA',
'EnerJoy', 'Lemon', 'Zesty and uplifting with a hint of joy', 'EnerJoy Nutritionals', 'Canada',
'BuzzBlend', 'Coffee', 'Rich and robust for a caffeine buzz', 'BuzzBlend Co.', 'USA',
'ZingZang', 'Ginger', 'Spicy zing with a tangy twist', 'ZingZang Innovations', 'Australia',
'PepPepper', 'Pepper', 'Spicy and fiery to pep you up', 'PepPepper Solutions', 'UK',
'MegaMint', 'Mint', 'Refreshingly cool and mega rejuvenating', 'MegaMint Beverages', 'USA',
'ChargeCherry', 'Cherry', 'Sweet and tart charge of fruity energy', 'ChargeCherry Ltd.', 'Canada',
'SparkSplash', 'Watermelon', 'Bursting with refreshing spark and splash', 'SparkSplash Co.', 'Australia',
'KickKola', 'Cola', 'Classic cola kick with an energy boost', 'KickKola Inc.', 'UK',
'ZapZen', 'Green Apple', 'Crisp and zen-inducing zap of energy', 'ZapZen Wellness', 'USA',
'VitaliTea', 'Herbal', 'Naturally invigorating herbal infusion', 'VitaliTea Blends', 'Canada',
'ThrillTonic', 'Lemon', 'Tart and thrilling tonic for a zestful day', 'ThrillTonic Ltd.', 'Australia',
'EnerBrew', 'Cocoa', 'Rich and chocolaty brew for sustained energy', 'EnerBrew Corp.', 'USA',
'BuzzBerry', 'Blueberry', 'Buzzing with antioxidant-rich berry power', 'BuzzBerry Beverages', 'Canada',
'PowerPop', 'Peach', 'Pop of peachy energy to power through the day', 'PowerPop Co.', 'UK',
'ReviveRush', 'Raspberry', 'Rush of raspberry delight to revive your senses', 'ReviveRush Inc.', 'Australia',
'ZapZing', 'Mango', 'Zingy mango explosion for instant rejuvenation', 'ZapZing Nutrition', 'USA',
'EnerGlow', 'Orange', 'Glowing with citrus energy to light up your day', 'EnerGlow Drinks', 'Canada',
'PunchPomegranate', 'Pomegranate', 'Punchy pomegranate flavor bursting with vitality', 'PunchPomegranate Ltd.', 'Australia',
'MintyMango', 'Mango', 'Minty mango fusion for a refreshing boost', 'MintyMango Corp.', 'USA',
'ZestZen', 'Lime', 'Zesty lime essence infused with zen vibes', 'ZestZen Beverages', 'Canada',
'EnerFresh', 'Grapefruit', 'Fresh and tangy grapefruit zest to revitalize', 'EnerFresh Co.', 'UK',
'PowerPear', 'Pear', 'Powerful pear punch to energize your day', 'PowerPear Nutrition', 'USA',
'Zingology', 'Kiwi', 'Zingy kiwi delight packed with energizing nutrition', 'Zingology Inc.', 'Canada',
'EnerChill', 'Cucumber', 'Chill with a burst of cool cucumber energy', 'EnerChill Wellness', 'Australia');

UPDATE compositions
SET energetics_id = (
    SELECT energetics_id
    FROM energetics
    WHERE energetics.composition_id = compositions.composition_id
);

select * from energetics
	

	
