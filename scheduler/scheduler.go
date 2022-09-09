package scheduler

import (
	"certcheckerbot/storage"
	"log"
	"time"
)

//Initialize the scheduler for every hour on the border of the next hour
func InitScheduler(db storage.UsersConfig, usersDomainsChan chan *storage.User) {
	now := time.Now()
	duration := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location()).Sub(now)
	//duration := time.Duration(time.Second) //string for tests
	log.Printf("Init scheduler after %v", duration)
	select {
	case <-time.After(duration):
		log.Println("Scheduler initialised!")
		startHourlyCheck(db, usersDomainsChan)
		for tick := range time.Tick(time.Hour) {
			log.Println("New scheduler tick " + tick.String())
			startHourlyCheck(db, usersDomainsChan)
		}
	}
}

func startHourlyCheck(db storage.UsersConfig, usersDomainsChan chan *storage.User) {
	location, err := time.LoadLocation("")
	if err != nil {
		log.Println(err)
		return
	}

	checkedUsersDomains := getCheckedUsersDomains(db, time.Now().In(location).Hour())

	if checkedUsersDomains != nil {
		for _, user := range *checkedUsersDomains {
			if user.UserDomains != nil {
				log.Println("Check domains for user " + user.Name)
				usersDomainsChan <- &user
			}
		}
	}
}

func getCheckedUsersDomains(db storage.UsersConfig, utcHour int) *[]storage.User {
	schedules, err := db.GetUsersSchedules()
	if err != nil {
		if err != storage.ErrorUsersSchedulesNotFound {
			log.Println(err)
		}
		return nil
	}

	var checkedUsersDomains []storage.User

	if schedules != nil {
		for _, schedule := range *schedules {
			if (schedule.NotificationHour - schedule.UTC) == utcHour {
				user, err := db.GetUserById(schedule.UserId)
				if err != nil {
					log.Println(err)
					continue
				}
				userDomains, err := db.GetUserDomains(user)
				if err != nil {
					log.Println(err)
					continue
				}
				user.UserDomains = *userDomains
				checkedUsersDomains = append(checkedUsersDomains, *user)
			}
		}
		return &checkedUsersDomains
	}

	return nil
}
