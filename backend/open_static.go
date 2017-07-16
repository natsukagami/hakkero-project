package backend

import (
	"math/rand"
	"time"
)

// This file contains a static list of open sentences, and a provider that returns a random one.

type staticOpenSentence struct {
	rnd rand.Source
}

func (s *staticOpenSentence) OpenSentence() (string, error) {
	n := s.rnd.Int63() % int64(len(staticOPs))
	return staticOPs[n], nil
}

// StaticOP returns an OpenSentencer that grabs static open sentences.
func StaticOP() OpenSentencer {
	rnd := rand.NewSource(time.Now().Unix())
	return &staticOpenSentence{rnd}
}

var staticOPs = []string{
	"I walked to work today.",
	"The ringing phone filled her with dread.",
	"I knew what I'd done as soon as the door closed.",
	"I remember the day.",
	"They turned and hurried back down the steps.",
	"The birds swoop low.",
	"He didn't look anything like she expected.",
	"The scent of lavender was overpowering.",
	"It had to be done.",
	"Carol enjoyed playing practical jokes.",
	"It wasn't as if anyone got hurt.",
	"Marg threw her bag on the floor and burst into tears.",
	"The house dwarfed everything in the street.",
	"It was never going to be an ordinary day.",
	"The sound of breaking glass stopped her.",
	"This wasn't where he wanted to be.",
	"Her foot slipped and she started to fall.",
	"The lights appeared out of the darkness.",
	"The clock struck one.",
	"The water looked deep and inviting.",
	"She stared into the darkness.",
	"In the end it didn't matter.",
	"The old woman turned and smiled.",
	"Everyone in the office turned and stared.",
	"He laughed in my face.",
	"I couldn't see any other way out of this mess.",
	"Mark didn't let on that he was scared.",
	"I can't help it. I lie. All the time.",
	"When I found out there would be a supermoon in two nights, I began making my plans.",
	"I turned to see who was following me.",
	"I've never done anything like this, but I was about to be thrown out of college and was desperate for money…",
	"The day I decided to get my tattoo…",
	"I had always thought the people who were paid to watch me were stupid, but this was beyond belief.",
	"I wish I could take back that moment, at the fortune-teller's table…",
	"Anytime you want to meet someone over the internet, take my advice. Don't.",
	"It's bad enough that I have the boss from hell, but this profession sucks the life from my soul, especially today.",
	"The first thing that went through my head was \"she's a witch!\" from that Monty Python movie.",
	"Manipulating people is so easy I almost stopped doing it. Almost.",
	"This was the last thing I expected.",
	"When he suggested I should run for office, I laughed. Then I considered it. After all…",
	"I was fascinated by the history of my house.",
	"This building smelled like a hospital, but it was not a place anyone would leave alive.",
	"I'd always wondered what would happen when I opened that door.",
	"The day I died began as an ordinary day.",
	"The moon has always called to me…",
	"Everyone thinks I'm normal, but no one has ever seen me at midnight.",
	"As I sit down to write this, I imagine what you will think when you open the envelope…",
	"Watching for the delivery truck became my obsession. I couldn't wait for the package to come; I'd never ordered anything like it before.",
	"Never a dull moment when you're a taxi driver. Just the other day this guy gets in…",
	"This might seem like it's about me, but it's not.",
	"If I hadn't looked out the window at that exact moment and watched it happen with my own eyes, I would never have believed it.",
	"Sometimes it's best not to go home again.",
	"I wasn't allowed to do this. But I couldn't resist.",
	"I never should have started it, but I had no idea I would cause such trouble.",
	"When I woke, I didn't know where I was.",
	"The moment I realized what I was reading, I knew I was as good as dead.",
	"I'd never imagined I could kill somebody.",
	"After she'd told me it was \"high time\" I knew my family legacy, my grandmother turned, pulled the box from the closet and handed it to me.",
}
