import PlusIcon from "../icons/PlusIcon.tsx";
import {useMemo, useState} from "react";
import {Column, Id, Task} from "../types.ts";
import ColumnContainer from "./ColumnContainer.tsx";
import axios from 'axios';


import {
    DndContext,
    DragEndEvent, DragOverEvent,
    DragOverlay,
    DragStartEvent,
    PointerSensor,
    useSensor,
    useSensors
} from "@dnd-kit/core";
import {arrayMove, SortableContext} from "@dnd-kit/sortable";
import {createPortal} from "react-dom";
import TaskCard from "./TaskCard.tsx";

function KanbanBoard() {
    const [columns, setColumns] = useState<Column[]>([]);
    const columnsId = useMemo(() => columns.map(col => col.id), [columns]);

    const [tasks, setTasks] = useState<Task[]>([]);

    const [activeColumn, setActiveColumn] = useState<Column | null>(null);

    const [activeTask, setActiveTask] = useState<Task | null>(null);

    const sensors = useSensors(
        useSensor(PointerSensor, {
            activationConstraint: {
                distance: 3 // 3px before move
         }
        }));

    return (
        <div
            className="
                m-auto
                flex
                min-h-screen
                w-full
                items-center
                overflow-x-auto
                overflow-y-hidden
                px-[40px]
            "
        >
            <DndContext onDragStart={onDragonStart} onDragEnd={onDragonEnd} onDragOver={onDragOver} sensors={sensors}>
                <div className="m-auto flex gap-4">
                    <div className="flex gap-4">
                        <SortableContext items = {columnsId}>
                            {columns.map((col) => (
                                <ColumnContainer
                                    key={col.id}
                                    column={col}
                                    deleteColumn={deleteColumn}
                                    updateColumn={updateColumn}
                                    createTask={createNewTask}
                                    deleteTask={deleteTask}
                                    updateTask={updateTask}
                                    tasks={tasks.filter((task) => task.columnId === col.id)}
                                />
                            ))}
                        </SortableContext>
                    </div>
                    <button
                        onClick={() => {
                            createNewColumn();
                        }}
                        className="
                        h-[60px]
                        w-[360px]
                        min-w-[350px]
                        cursor-pointer
                        rounded-lg
                        bg-mainBackgroundColor
                        border-2
                        border-columnBackgroundColor
                        p-4
                        ring-rose-500
                        hover:ring-2
                        flex
                        gap-2
                        "
                    >
                        <PlusIcon />
                        Добавить столбец
                    </button>
                </div>
                {createPortal(
                    <DragOverlay>
                        {activeColumn && (
                            <ColumnContainer
                                column={activeColumn}
                                deleteColumn={deleteColumn}
                                updateColumn={updateColumn}
                                createTask={createNewTask}
                                deleteTask={deleteTask}
                                updateTask={updateTask}
                                tasks={tasks.filter((task) => task.columnId === activeColumn.id)}
                            />
                        )}
                        {activeTask &&
                            <TaskCard
                                task={activeTask}
                                deleteTask={deleteTask}
                                updateTask={updateTask}
                            />
                        }
                    </DragOverlay>,
                    document.body
                )}
            </DndContext>
        </div>
    );

    function createNewTask(columnId: Id) {
        const  description = `Задача ${tasks.length + 1}`
        const requestData = {
            column_id: columnId,
            description: description,
        };
        console.log(123)
        axios.post('http://localhost:5007/api/v1/task/create', requestData)
            .then((response) => {
                if (response.status === 201) {
                    const newTask: Task = {
                        id: response.data.id,
                        columnId,
                        content: description
                    }
                    setTasks([...tasks, newTask]);
                } else {
                    console.error('Неправильный статус ответа:', response.status);
                }
            })
            .catch((error) => {
                console.error('Ошибка при отправке запроса:', error);
            });
    }

    function deleteTask(id: Id) {
        axios.delete(`http://localhost:5007/api/v1/task/${id}`)
            .then((response) => {
                if (response.status === 200) {
                    const newTasks = tasks.filter(task => task.id !== id);
                    setTasks(newTasks);
                } else {
                    console.error('Неправильный статус ответа:', response.status);
                }
            })
            .catch((error) => {
                console.error('Ошибка при отправке запроса:', error);
            });
    }

    function updateTask(id: Id, content: string) {
        const requestData = {
            description: content,
        };
        axios.patch(`http://localhost:5007/api/v1/task/${id}/update_description`, requestData)
            .then((response) => {
                if (response.status === 200) {
                    const newTasks = tasks.map(task => {
                        if (task.id !== id){
                            return task;
                        }
                        return {...task, content};
                    });
                    setTasks(newTasks);
                } else {
                    console.error('Неправильный статус ответа:', response.status);
                }
            })
            .catch((error) => {
                console.error('Ошибка при отправке запроса:', error);
            });

        const newTasks = tasks.map(task => {
            if (task.id !== id){
                return task;
            }
            return {...task, content};
        });

        setTasks(newTasks);
    }

    function createNewColumn(){ //TODO принимать user_id...или сделать user_id глобально
        const requestData = {
            user_id: 1, //TODO hadrcode!!!
        };
        axios.post('http://localhost:5007/api/v1/column/create', requestData)
            .then((response) => {
                if (response.status === 201) {
                    const newColumn = {
                        id: response.data.id,
                        title: `Столбец ${columns.length + 1}`,
                    };
                    setColumns([...columns, newColumn]);
                } else {
                    console.error('Неправильный статус ответа:', response.status);
                }
            })
            .catch((error) => {
                console.error('Ошибка при отправке запроса:', error);
            });
    }

    function deleteColumn(id: Id){
        axios.delete(`http://localhost:5007/api/v1/column/${id}`)
            .then((response) => {
                if (response.status === 200) {
                    const filteredColumns = columns.filter(col => col.id !== id);
                    setColumns(filteredColumns);

                    const newTasks = tasks.filter(t => t.columnId !== id);
                    setTasks(newTasks);
                } else {
                    console.error('Неправильный статус ответа:', response.status);
                }
            })
            .catch((error) => {
                console.error('Ошибка при отправке запроса:', error);
            });
    }


    function updateColumn(id:Id, title: string) {
        const requestData = {
            user_id: 1, //TODO hadrcode!!!
            name: title,
        };
        axios.patch(`http://localhost:5007/api/v1/column/${id}/update_name`, requestData)
            .then((response) => {
                if (response.status === 200) {
                    const newColumns = columns.map((col) => {
                        if (col.id !== id) return col;
                        return {...col, title };
                    });
                    setColumns(newColumns);
                } else {
                    console.error('Неправильный статус ответа:', response.status);
                }
            })
            .catch((error) => {
                console.error('Ошибка при отправке запроса:', error);
            });
    }

    function onDragonStart(event: DragStartEvent){
        console.log("DRAG START", event);
        if (event.active.data.current?.type === "Column"){
            setActiveColumn(event.active.data.current.column);
            return;
        }

        if (event.active.data.current?.type === "Task"){
            setActiveTask(event.active.data.current.task);
            return;
        }
    }

    function onDragonEnd(event: DragEndEvent){
        setActiveColumn(null);
        setActiveTask(null);

        const {active, over} = event;

        if (!over) return;

        const activeId = active.id;
        const overId = over.id;

        if (activeId === overId){
            return;
        }

        setColumns(columns => {
            const activeColumnIndex = columns.findIndex((col) => col.id === activeId);
            const overColumnIndex = columns.findIndex((col) => col.id === overId);

            return arrayMove(columns, activeColumnIndex, overColumnIndex);
        });
    }

    function onDragOver(event: DragOverEvent) {
        const {active, over} = event;

        if (!over) return;

        const activeId = active.id;
        const overId = over.id;

        if (activeId === overId){
            return;
        }

        const isActiveATask = active.data.current?.type === "Task";
        const isOverATask = over.data.current?.type === "Task";

        if (!isActiveATask) {
            return;
        }

        //1)  Dropping a Task over another Task
        if (isActiveATask && isOverATask) {
            setTasks(tasks => {
                const activeIndex = tasks.findIndex((t) => t.id === activeId);
                const overIndex = tasks.findIndex((t) => t.id === overId);

                tasks[activeIndex].columnId = tasks[overIndex].columnId;

                return arrayMove(tasks, activeIndex, overIndex);
            });
        }

        const isOverAColumn = over.data.current?.type === "Column";

        //2)  Dropping a Task over a column
        if (isActiveATask && isOverAColumn) {
            const requestData = {
                column_id: overId,
            };
            axios.patch(`http://localhost:5007/api/v1/task/${activeId}/update_column_id`, requestData)
                .then((response) => {
                    if (response.status !== 200) {
                        console.error('Неправильный статус ответа:', response.status);
                    }
                })
                .catch((error) => {
                    console.error('Ошибка при отправке запроса:', error);
                });
            setTasks(tasks => {
                const activeIndex = tasks.findIndex((t) => t.id === activeId);

                tasks[activeIndex].columnId = overId;

                return arrayMove(tasks, activeIndex, activeIndex);
            });
        }

    }
}

export default KanbanBoard