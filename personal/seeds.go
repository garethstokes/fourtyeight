package personal

func garrydanger() * Person {
  person := new( Person )
  person.Username = "@garrydanger"
  person.Email = "garrydanger@gmail.com"
  person.AvatarUrl = "https://si0.twimg.com/profile_images/2083020030/Photo_on_2012-03-16_at_15.47__2.jpg"
  return person
}

func shredder() * Person {
  person := new( Person )
  person.Username = "@shredder"
  person.Email = "shredder@gmail.com"
  person.AvatarUrl = "https://si0.twimg.com/profile_images/1434628104/zzzz-_3_.png"
  return person
}

func caveman() * Person {
  person := new( Person )
  person.Username = "@caveman"
  person.Email = "big_scary_cave@gmail.com"
  person.AvatarUrl = "https://trello-avatars.s3.amazonaws.com/dd2ab4b70b525b89fa68abf63d259d7e/original.png"
  return person
}

func (s * Personal) seedSingleUser( user * Person ) {
  _, error := s.Create( user, "bobafett" )
  if error != nil {
    s.logf( "Creating user, %s... ERROR\n", user.Username )
    s.logf( "%s\n", error )
    return
  }

  s.logf( "Creating user, %s... SUCCESS\n", user.Username )
}

func (s * Personal) Seed() {
  s.log( "\nSeeding Personal" )
  s.log( "================" )

  s.seedSingleUser( garrydanger() )
  s.seedSingleUser( shredder() )
  s.seedSingleUser( caveman() )

  s.log( "\n\n" )
}
