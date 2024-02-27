package models

import (
	"gorm.io/gorm"
	"job_board/db"
)

var database *gorm.DB

func init() {
	database = db.GetDB()
}

func MigrateDb() {

	//remeber to run this for uuid support CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	database.AutoMigrate(
		&User{},

		&Company{},
		&Industry{},
		&EmployeesSize{},
		&JobType{},
		&Level{},
		&Job{},
		&JobApplication{},

		&Profile{},
		&SalaryCurrency{},
		&Gender{},
		&Degree{},
		&AcademicRanking{},
		&Education{},
		&InternShipExperience{},
		&ProjectsExperience{},
		&WorkSample{},
		&Award{},
		&Language{},
		&LanguageProficiency{},
		&ProfileLanguage{},
		&SocialMedia{},
		&SocialMediaAccount{})
		
}
