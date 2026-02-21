export interface Exercise {
  id: number
  title: string
  description: string
}

export interface Segment {
  id: number
  title: string
  exercise_id: number
  spoken_text: string
  audio_file: string
}

export interface SegmentAnswer {
  exercise_id: number
  segment_id: number
  transcript: string
  rating: number
}
