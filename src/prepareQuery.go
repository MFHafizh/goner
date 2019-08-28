package src

var queries map[string]string

func prepareQueryString() {
	queries = map[string]string{
		"InsertLabel": `
			INSERT INTO activity_label
				(activity_id, category_id, category_type) 
			VALUES
				($1, $2, $3)
			ON CONFLICT ON CONSTRAINT activity_label_pkey DO UPDATE SET category_id = $2, category_type = $3
		`,
		"RemoveLabelAct": `
			DELETE FROM activity_label
			WHERE activity_id = $1
		`,
		"GetCategoriesLabel": `
			SELECT
				category_id,
				category_type
			FROM activity_label	
			WHERE activity_id = $1
		`,
		"GetAllActivityLabel": `
			SELECT
				ac.id AS actID,
				al.category_id AS categoryID,
				al.category_type AS categoryType,
				ac.create_time AS showTime
			FROM activities ac
			INNER JOIN activity_label al
				ON (al.activity_id = ac.id)
			WHERE
				ac.status = 1 AND
				ac.content_type_id = 13
			UNION
			SELECT 
				act.id AS actID,
				COALESCE(al.category_id,0) AS categoryID,
				COALESCE(al.category_type,0) as categoryType,
				cp.show_time AS showTime
			FROM activities act 
			LEFT JOIN activity_label al
				ON al.activity_id = act.id
			INNER JOIN content_post cp
				ON cp.id = ANY(act.content::INT[]) 
			WHERE 
				act.status = 1 AND
				act.content_type_id IN (20,25) AND
				cp.status = 1 AND
				cp.show_in_explore = true
			ORDER BY showTime DESC , actID DESC
		`,
		"GetContentDetail": `
			SELECT content FROM activity_custom WHERE id = $1
		`,
		"GetAllActivitiesContent": `
			SELECT ac.id
			FROM activities ac
			INNER JOIN kol_activity_qc qc ON ac.id = qc.id_activity
		`,
		"GetDetailActivitiesContent": `
			SELECT ac.content
			FROM activities ac
			INNER JOIN kol_activity_qc qc ON ac.id = qc.id_activity
			WHERE ac.id = $1
		`,
		"GetDetailActivityCustom": `
			SELECT
				ac.create_by AS create_by,
				ac.create_time AS create_time,
				ac.update_by AS update_by,
				ac.update_time AS update_time
			FROM
				activity_custom ac
			WHERE ac.id = $1
		`,
		"GetDetailLabelActivityQC": `
			SELECT category_id, category_type FROM activity_label al
			INNER JOIN kol_activity_qc qc ON qc.id_activity = al.activity_id
			WHERE qc.id = $1;
		`,
		"GetAllContentNotLabelled": `
			SELECT qc.id_activity FROM  kol_activity_qc qc
			WHERE qc.id_activity NOT IN (SELECT activity_id FROM activity_label al);
		`,
		"GetLabelFromActivityID": `
			SELECT 
				al.category_id, al.category_type
			FROM
				activity_label al
			WHERE
				al.activity_id = $1;
		`,
	}
}
