export interface Segment {
  id: number
  title: string
  exercise_id: number
  spoken_text: string
  audio_file: string
}

export interface SegmentAnswer {
  rating: number
  transcript: string
}
