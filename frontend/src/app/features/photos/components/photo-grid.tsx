interface PhotoGridProps {
  images: string[];
}

export function PhotoGrid({ images }: PhotoGridProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
      {images.map((src, i) => (
        <img key={i} src={src} alt={`Photo ${i + 1}`} className="w-full rounded-md" />
      ))}
    </div>
  );
}
