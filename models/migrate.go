package models

import (
	"job_board/db"

	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	database = db.GetDB()
}

func MigrateDb() {

	//remeber to run this for uuid support CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	//this is for index
	// if err := database.Exec("CREATE INDEX idx_user_email ON users (email) WHERE email IS NOT NULL").Error; err != nil {
	// 	log.Printf("Error creating index on users table: %v", err)
	// 	return
	// }

	// result := database.Unscoped().Where("role_name = ?", SuperAdminRole).Delete(&User{})
	// // Check for errors during delete operation
	// if err := result.Error; err != nil {
	// 	log.Printf("Error deleting admin: %v", err)
	// 	return
	// }
	database.AutoMigrate(
	// &User{},

	// &Company{},
	// &Industry{},
	// &EmployeesSize{},
	// &JobType{},
	// &Level{},
	// &Job{},
	// &JobApplication{},

	// &Profile{},
	// &SalaryCurrency{},
	// &Gender{},
	// &Degree{},
	// &AcademicRanking{},
	// &Education{},
	// &InternShipExperience{},
	// &ProjectsExperience{},
	// &WorkSample{},
	// &Award{},
	// &Language{},
	// &LanguageProficiency{},
	// &ProfileLanguage{},
	// &SocialMedia{},
	// &SocialMediaAccount{},
	)

}
