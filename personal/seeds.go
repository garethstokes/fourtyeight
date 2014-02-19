package personal

func make_garrydanger() * Person {
  person := new( Person )
  person.Username = "garrydanger"
  person.Email = "garrydanger@gmail.com"
  person.AvatarUrl = "https://si0.twimg.com/profile_images/2083020030/Photo_on_2012-03-16_at_15.47__2.jpg"

  person.Followers = []string {"shredder", "caveman"}
  person.Following = []string {"shredder", "caveman"}

  return person
}

func make_shredder() * Person {
  person := new( Person )
  person.Username = "shredder"
  person.Email = "shredder@gmail.com"
  person.AvatarUrl = "https://si0.twimg.com/profile_images/1434628104/zzzz-_3_.png"

  person.Followers = []string {"garrydanger", "caveman"}
  person.Following = []string {"garrydanger", "caveman"}

  return person
}

func make_caveman() * Person {
  person := new( Person )
  person.Username = "caveman"
  person.Email = "big_scary_cave@gmail.com"
  person.AvatarUrl = "https://trello-avatars.s3.amazonaws.com/dd2ab4b70b525b89fa68abf63d259d7e/original.png"

  person.Followers = []string {"shredder", "garrydanger"}
  person.Following = []string {"shredder", "garrydanger"}

  return person
}

func (s * Personal) seedSingleUser( user * Person ) * Person {
  person, error := s.Create( user, "bobafett" )
  if error != nil {
    s.logf( "Creating user, %s... ERROR\n", user.Username )
    s.logf( "%s\n", error )
    return nil
  }

  s.logf( "Creating user, %s... SUCCESS\n", user.Username )
  return person
}

func (s * Personal) Seed() {
  s.log( "\nSeeding Personal" )
  s.log( "================" )

  s.seedSingleUser( make_garrydanger() )
  s.seedSingleUser( make_shredder() )
  s.seedSingleUser( make_caveman() )

  s.log( "\n\n" )
}
