package course_material

import "context"

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM course_materials WHERE id = $1", id)
	return err
}
